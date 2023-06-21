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

var clients = make(map[*websocket.Conn]bool) // Connected clients

var upgrader = websocket.Upgrader{} // WebSocket upgrader

func NewUserInformation(r *mux.Router) {
	r.HandleFunc("/test/{text}", testHandler).Methods("GET")

	r.HandleFunc("/ws", testWebsocket)
}

func testHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	text := params["text"]

	body := application.NewResponse(fmt.Sprintf("test: %v", text))

	writeReponse(w, r, body)
}

func testWebsocket(w http.ResponseWriter, r *http.Request) {

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	// Register client
	clients[conn] = true

	// Close connection and remove client when done
	defer func() {
		delete(clients, conn)
		conn.Close()
	}()

	// Create a file to write the messages
	file, err := os.Create("messages.log")
	if err != nil {
		log.Println("File creation error:", err)
		return
	}
	defer file.Close()

	// Read messages from client
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		// // Broadcast the received message to all connected clients
		// broadcast <- msg
		// Write the message to the file
		_, err = file.WriteString(time.Now().Format(time.RFC3339) + ": " + string(msg) + "\n")
		if err != nil {
			log.Println("File write error:", err)
			break
		}
	}

}
