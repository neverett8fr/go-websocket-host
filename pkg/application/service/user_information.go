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
		log.Println("WebSocket connection not found")
		return
	}

	conn := ww.conn
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println("WebSocket read error:", err)
			}
			break
		}

		log.Println("Received message:", string(message))

		// Additional logic to process the received message
		// ...
	}
}
