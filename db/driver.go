package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var uri = "mongodb://localhost:27017/?maxPoolSize=20&w=majority"

type DB struct {
	Client  *mongo.Client
	Context context.Context
}

var dbConn = &DB{
	Client:  nil,
	Context: context.TODO(),
}

func ConnectToClient() (*DB, error) {

	client, err := mongo.Connect(dbConn.Context, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to db client...")

	dbConn.Client = client

	return dbConn, err
}

func PingBD(db *DB) {
	if err := db.Client.Ping(dbConn.Context, readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Client ping!!!...")
}

func SetDbUri(user, password, host string) {
	uri = fmt.Sprintf("mongodb://%s:%s@%s:27017/authentication?maxPoolSize=20&w=majority", user, password, host)
}
