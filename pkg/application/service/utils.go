package service

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"salve-data-service/pkg/infra/db"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	DBConn *db.DBConn

	WebsocketServer *WebSocketServer
	Upgrader        websocket.Upgrader
)

func NewServiceRoutes(r *mux.Router, conn *sql.DB) {

	DBConn = db.NewDBConnFromExisting(conn)
	WebsocketServer = NewWebsocketServer()

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
