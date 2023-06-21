package service

import (
	"log"
	"net/http"

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

func WebsocketMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Perform middleware actions before upgrading to WebSocket

		// Example: Log the incoming request
		log.Println("Incoming request:", r.URL.Path)

		// Upgrade HTTP connection to WebSocket
		conn, err := Websockets.Upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("WebSocket upgrade error:", err)
			return
		}

		defer func() {
			// Perform middleware actions after serving the WebSocket connection

			// Example: Log WebSocket connection closure
			conn.Close()
			log.Println("WebSocket connection closed")
		}()

		// Perform additional actions after upgrading to WebSocket, if needed
		// Register client
		Websockets.Clients[conn] = true

		// Example: Log successful WebSocket upgrade
		log.Println("WebSocket upgrade successful")

		// Create a new request with the upgraded WebSocket connection
		req := r.WithContext(r.Context())
		req.Header.Set("Upgrade", "websocket")
		req.Header.Set("Connection", "Upgrade")

		// Create a WebSocketResponseWriter
		ww := &WebSocketResponseWriter{
			ResponseWriter: w,
			conn:           conn,
		}

		// Serve the WebSocket connection
		next.ServeHTTP(ww, req)
	})
}
