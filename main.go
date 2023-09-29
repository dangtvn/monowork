package main

import (
	"encoding/json"
	"log"
	"monowork/internal"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func main() {
	router := mux.NewRouter()
	var ServerPort = os.Getenv("PORT")

	if ServerPort == "" {
		ServerPort = "4444"
	}

	go internal.MusicStation.Run()

	router.HandleFunc("/new-playing", currentTrackHandler)
	router.HandleFunc("/stream.mp3", streamHandler)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./www/dist/")))

	log.Printf("The server is streaming on http://localhost:%s", ServerPort)
	log.Fatal(http.ListenAndServe(":"+ServerPort, router), nil)
}

func currentTrackHandler(w http.ResponseWriter, r *http.Request) {
	trackInfo, _ := internal.MusicStation.TrackInfo()
	resp, _ := json.Marshal(trackInfo.Duration())
	w.Header().Set("Content-Type", "application/json")

	w.Write(resp)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	client := internal.NewClient(ws)
	internal.MusicStation.Register(client)
	go client.ReadPump()
	client.WritePump()
}
