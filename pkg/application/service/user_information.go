package service

import (
	"fmt"
	"log"
	"net/http"
	"os"
	application "salve-data-service/pkg/application/entities"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func newUserInformation(r *mux.Router) {
	r.HandleFunc("/test/{text}", testHandler).Methods("GET")

	subRouter := r.PathPrefix("/ws").Subrouter()
	subRouter.Use(WebsocketMiddleware)
	subRouter.HandleFunc("/test", testWebsocket)
	subRouter.HandleFunc("/file", testWriteToFile)
	subRouter.HandleFunc("/db", testWriteToDB)
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

func testWriteToFile(w http.ResponseWriter, r *http.Request) {
	ww, ok := w.(*WebSocketResponseWriter)
	if !ok {
		log.Println("Websocket connection not found")
		return
	}

	// Create a file to write the messages
	file, err := os.Create("messages.log")
	if err != nil {
		log.Println("File creation error:", err)
		return
	}
	defer file.Close()

	for {
		_, message, err := ww.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println("WebSocket read error:", err)
			}
			break
		}

		log.Println("Received message:", string(message))
		// Write the message to the file
		_, err = file.WriteString(time.Now().Format(time.RFC3339) + ": " + string(message) + "\n")
		if err != nil {
			log.Println("File write error:", err)
			break
		}

	}

}

func testWriteToDB(w http.ResponseWriter, r *http.Request) {
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
		err = DBConn.CreateEntry(string(message))
		if err != nil {
			log.Println(err)
			break
		}

	}

}
