package internal

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	ws   *websocket.Conn
	send chan []byte
}

func NewClient(ws *websocket.Conn) *Client {
	return &Client{
		ws:   ws,
		send: make(chan []byte, 256),
	}
}

func (c *Client) removeClient() {
	MusicStation.unregister <- c

	c.ws.Close()
}

func (c *Client) ReadPump() {
	defer c.removeClient()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// listen for socket
	for {
		_, _, err := c.ws.ReadMessage()
		if err != nil {
			log.Print(err.Error())
			break
		}
	}
}

func (c *Client) write(messageType int, message []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(messageType, message)
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	// write to socket
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			err := c.write(websocket.BinaryMessage, message)
			if err != nil {
				return
			}

		case <-ticker.C:
			err := c.write(websocket.PingMessage, []byte{})
			if err != nil {
				return
			}
		}
	}
}
