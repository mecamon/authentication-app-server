package main

import (
	"flag"
	"github.com/authentication-app-server/api-services/auth"
	"github.com/authentication-app-server/config"
	"github.com/authentication-app-server/db"
	"github.com/authentication-app-server/services"
	"log"
	"net/http"
	"os"
)

func main() {

	var (
		cloud      string
		cloudKey   string
		secret     string
		dbHost     string
		dbUser     string
		dbPassword string
		production bool
	)

	flag.StringVar(&cloud, "cloud", "", "cloud name")
	flag.StringVar(&cloudKey, "cloud-key", "", "cloud API key")
	flag.StringVar(&secret, "secret", "", "cloud API secret")
	flag.StringVar(&auth.ClientID, "client-id", "", "cloud app id")
	flag.StringVar(&auth.ClientSecret, "client-secret", "", "Client secret key")
	flag.StringVar(&dbHost, "db-host", "localhost", "Database host")
	flag.StringVar(&dbUser, "db-user", "", "Db username")
	flag.StringVar(&dbPassword, "db-password", "", "Db password")
	flag.BoolVar(&production, "production", true, "Environment mode")
	flag.Parse()

	if flag.Lookup("test.v") == nil {
		config.AppConfig.Production = production
		db.SetDbUri(dbUser, dbPassword, dbHost)
	}

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Running server in port", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal(err)
	}
}
