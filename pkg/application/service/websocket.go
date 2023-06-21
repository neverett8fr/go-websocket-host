package service

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketResponseWriter is a custom http.ResponseWriter wrapper
type WebSocketResponseWriter struct {
	http.ResponseWriter
	conn *websocket.Conn
}

// Write implements the Write method of http.ResponseWriter
func (w *WebSocketResponseWriter) Write(p []byte) (int, error) {
	err := w.conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

type Client struct {
	conn *websocket.Conn
}
type WebSocketServer struct {
	clients     map[*Client]bool
	clientsLock sync.Mutex
}

func NewWebsocketServer() *WebSocketServer {
	return &WebSocketServer{
		clients: make(map[*Client]bool),
	}
}

// addClient adds a new client to the WebSocket server
func (ws *WebSocketServer) addClient(client *Client) {
	ws.clientsLock.Lock()
	defer ws.clientsLock.Unlock()

	ws.clients[client] = true
	log.Println("Client connected. Total clients:", len(ws.clients))
}

// removeClient removes a client from the WebSocket server
func (ws *WebSocketServer) removeClient(client *Client) {
	ws.clientsLock.Lock()
	defer ws.clientsLock.Unlock()

	delete(ws.clients, client)
	log.Println("Client disconnected. Total clients:", len(ws.clients))
}

// broadcastMessage sends a message to all connected clients
func (ws *WebSocketServer) broadcastMessage(message []byte) {
	ws.clientsLock.Lock()
	defer ws.clientsLock.Unlock()

	for client := range ws.clients {
		err := client.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Error broadcasting message to client:", err)
			continue
		}
	}
}
