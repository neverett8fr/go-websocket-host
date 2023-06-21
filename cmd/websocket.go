package cmd

import (
	"github.com/gorilla/websocket"
)

func StartWebsocket() (map[*websocket.Conn]bool, websocket.Upgrader, error) {
	clients := make(map[*websocket.Conn]bool) // Connected clients
	upgrader := websocket.Upgrader{}          // WebSocket upgrader

	return clients, upgrader, nil
}
