package main

import (
	"github.com/authentication-app-server/db"
	"log"
	"net/http"
)

const port = ":8080"

func main() {

	dbConn, err := db.ConnectToClient()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := dbConn.Client.Disconnect(dbConn.Context); err != nil {
			panic(err)
		}
	}()

	//Testing the db connection
	db.PingBD(dbConn)

	router := makeRouter(dbConn)

	log.Println("Running server in port", port)
	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal(err)
	}
}
