package service

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Sockets struct {
	Clients  map[*websocket.Conn]bool
	Upgrader websocket.Upgrader
}

var (
	Websockets Sockets
)

func NewServiceRoutes(r *mux.Router, sockClient map[*websocket.Conn]bool, sockUpgrader websocket.Upgrader) {

	Websockets.Clients = sockClient
	Websockets.Upgrader = sockUpgrader

	newUserInformation(r)
}

func writeReponse(w http.ResponseWriter, r *http.Request, body interface{}) {

	reponseBody, err := json.Marshal(body)
	if err != nil {
		log.Printf("error converting reponse to bytes, err %v", err)
	}
	w.Header().Add("Content-Type", "application/json")

	_, err = w.Write(reponseBody)
	if err != nil {
		log.Printf("error writing response, err %v", err)
	}
}
