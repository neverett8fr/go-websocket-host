package service

import (
	"log"
	"net/http"
)

func WebsocketMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Example: Log the incoming request
		log.Println("Incoming request:", r.URL.Path)

		conn, err := Upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("WebSocket upgrade error:", err)
			return
		}

		c := &Client{
			conn: conn,
		}
		WebsocketServer.addClient(c)
		defer func() {
			WebsocketServer.removeClient(c)
		}()

		// Create a new request with the upgraded WebSocket connection
		req := r.WithContext(r.Context())
		req.Header.Set("Upgrade", "websocket")
		req.Header.Set("Connection", "Upgrade")

		ww := &WebSocketResponseWriter{
			ResponseWriter: w,
			conn:           conn,
		}

		next.ServeHTTP(ww, req)
	})
}
