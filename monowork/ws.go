package monowork

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type WebSocket struct {
	Clients []*websocket.Conn
	Reader  Reader
	Logger  zap.Logger
}

func NewWebSocket(reader Reader, logger zap.Logger) *WebSocket {
	return &WebSocket{
		Reader: reader,
		Logger: logger,
	}
}

func (ws *WebSocket) SaveClient(conn *websocket.Conn) {
	ws.Clients = append(ws.Clients, conn)
}

func (ws *WebSocket) RemoveClient(conn *websocket.Conn) {
	var clients []*websocket.Conn
	for _, client := range ws.Clients {
		if client != conn {
			clients = append(clients, client)
		}
	}
	ws.Clients = clients
}

func (ws *WebSocket) SendDataToClient(conn *websocket.Conn) {
	musicInfoJson, err := json.Marshal(ws.Reader.GetSongInfo())
	if err != nil {
		ws.Logger.Error(err.Error())
	} else {
		err := conn.WriteMessage(websocket.TextMessage, musicInfoJson)
		if err != nil {
			ws.Logger.Error(err.Error())
		}
	}
}

func (ws *WebSocket) SendDataToAllClient() {
	for _, client := range ws.Clients {
		ws.SendDataToClient(client)
	}
}
