package internal

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/tcolgate/mp3"
)

const (
	skip = 0
)

type Station struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	track      string
	trackIndex int

	tracks []string
}

func newStation() Station {
	return Station{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		track:      "",
		trackIndex: 0,
		tracks:     readTracksAt("./music"),
	}
}

var MusicStation = newStation()

func readTracksAt(dir string) []string {
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	songs := make([]string, 0)
	for _, f := range files {
		songs = append(songs, dir+"/"+f.Name())
	}

	return songs
}

func (station *Station) Register(client *Client) {
	station.register <- client
}

func (station *Station) unregisterClient(c *Client) {
	if _, ok := station.clients[c]; ok {
		delete(station.clients, c)
		close(c.send)
	}
}

func (station *Station) isEmpty() bool {
	return len(station.clients) == 0
}

func (station *Station) count() int {
	return len(station.clients)
}

func (station *Station) SendMessage(message []byte) {
	if station.isEmpty() {
		return
	}

	station.broadcast <- message
}

func (station *Station) Run() {
	go station.stream()

	for {
		select {
		case c := <-station.register:
			station.clients[c] = true
		case c := <-station.unregister:
			station.unregisterClient(c)
		case m := <-station.broadcast:
			station.broadcastMessage(m)
		}
	}
}

func (station *Station) broadcastMessage(m []byte) {
	for c := range station.clients {
		select {
		case c.send <- m:
		default:
			close(c.send)
			delete(station.clients, c)
		}
	}
}

func (station *Station) stream() {
	for {
		if err := station.selectTrack(); err == nil {
			station.streamStep()
		}

		return
	}
}

func (station *Station) selectTrack() error {
	if station.trackIndex >= len(station.tracks) {
		log.Print("no more songs")
		return errors.New("no more songs")
	}

	station.track = station.tracks[station.trackIndex]
	station.trackIndex++

	return nil
}

func (station *Station) streamStep() {
	r, err := os.Open(station.track)
	defer r.Close()

	if err != nil {
		log.Println(err)
		return
	}

	d := mp3.NewDecoder(r)
	var f mp3.Frame
	skipped := 0

	for {
		if err := d.Decode(&f, &skipped); err != nil {
			log.Println(err)
			return
		}
		b := make([]byte, f.Size())

		f.Reader().Read(b)
		station.SendMessage(b)

		time.Sleep(f.Duration())
	}
}

func (station *Station) TrackInfo() (mp3.Frame, error) {
	var f mp3.Frame
	skipped := 0

	r, err := os.Open(station.track)
	if err != nil {
		log.Println(err)
		return f, err
	}

	d := mp3.NewDecoder(r)
	d.Decode(&f, &skipped)

	return f, nil
}
