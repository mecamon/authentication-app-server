package main

import (
	"flag"
	"github.com/authentication-app-server/api-services/auth"
	"github.com/authentication-app-server/db"
	"github.com/authentication-app-server/services"
	"log"
	"net/http"
)

const port = ":8080"

func main() {
	var (
		cloud    string
		cloudKey string
		secret   string
	)

	flag.StringVar(&cloud, "cloud", "", "cloud name")
	flag.StringVar(&cloudKey, "cloud-key", "", "cloud API key")
	flag.StringVar(&secret, "secret", "", "cloud API secret")
	flag.StringVar(&auth.ClientID, "client-id", "", "cloud app id")
	flag.StringVar(&auth.ClientSecret, "client-secret", "", "Client secret key")
	flag.Parse()

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

	err = services.NewCloudinaryInstance(cloud, cloudKey, secret)
	if err != nil {
		log.Println("Error creating cloudinary instance", err)
	}

	log.Println("Running server in port", port)
	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal(err)
	}
}
