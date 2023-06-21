package main

import (
	"log"
	"salve-data-service/cmd"
	application "salve-data-service/pkg/application/service"
	"salve-data-service/pkg/config"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Route declaration
func getRoutes(conn map[*websocket.Conn]bool, upg websocket.Upgrader) *mux.Router {
	r := mux.NewRouter()
	application.NewServiceRoutes(r, conn, upg)

	return r
}

// Initiate web server
func main() {
	conf, err := config.Initialise()
	if err != nil {
		log.Fatalf("error initialising config, err %v", err)
		return
	}
	log.Println("config initialised")

	serviceDB, err := cmd.OpenDB(&conf.DB)
	if err != nil {
		log.Fatalf("error starting db, err %v", err)
		return
	}
	defer serviceDB.Close()
	log.Println("connection to DB setup")

	err = cmd.MigrateDB(serviceDB, conf.DB.Driver)
	if err != nil {
		log.Fatalf("error running DB migrations, %v", err)
		return
	}
	log.Println("DB migrations ran")

	socketConn, socketUpgrader, err := cmd.StartWebsocket()
	if err != nil {
		log.Fatalf("error starting websocket server, err %v", err)
		return
	}

	router := getRoutes(socketConn, socketUpgrader)
	log.Println("API routes retrieved")

	err = cmd.StartServer(&conf.Service, router)
	if err != nil {
		log.Fatalf("error starting server, %v", err)
		return
	}
	log.Println("server started")

}
