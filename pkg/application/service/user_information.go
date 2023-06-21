package service

import (
	"fmt"
	"log"
	"net/http"
	application "salve-data-service/pkg/application/entities"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func newUserInformation(r *mux.Router) {
	r.HandleFunc("/test/{text}", testHandler).Methods("GET")

	subRouter := r.PathPrefix("/ws").Subrouter()
	subRouter.Use(WebsocketMiddleware)
	subRouter.HandleFunc("/test", testWebsocket)
}

func testHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	text := params["text"]

	body := application.NewResponse(fmt.Sprintf("test: %v", text))

	writeReponse(w, r, body)
}

func testWebsocket(w http.ResponseWriter, r *http.Request) {
	ww, ok := w.(*WebSocketResponseWriter)
	if !ok {
		log.Println("Websocket connection not found")
		return
	}

	for {
		_, message, err := ww.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println("WebSocket read error:", err)
			}
			break
		}

		log.Println("Received message:", string(message))

		// Broadcast the received message to all clients
		WebsocketServer.broadcastMessage(message)
	}
}
