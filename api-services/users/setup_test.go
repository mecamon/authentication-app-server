package users

import (
	"context"
	"github.com/authentication-app-server/db"
	"github.com/authentication-app-server/middlewares"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"testing"
	"time"
)

var mongoTestClient *mongo.Client
var mongoTestDBName = "authentication-test"
var testMainRouter *chi.Mux

func TestMain(m *testing.M) {

	dbConn, err := db.ConnectToClient()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := dbConn.Client.Disconnect(dbConn.Context); err != nil {
			panic(err)
		}
	}()
	mongoTestClient = dbConn.Client
	testMainRouter = createMainRouter(mongoTestClient, mongoTestDBName)

	code := m.Run()

	shutdown()

	os.Exit(code)
}

func createMainRouter(mongoClient *mongo.Client, dbName string) *chi.Mux {
	mainRouter := chi.NewRouter()
	router := chi.NewRouter()
	router.Use(middlewares.TokenValidation)

	NewHandlers(mongoClient, dbName)

	router.Put("/update", handlers.updateUserInfo)
	router.Get("/info", handlers.userInfo)
	mainRouter.Mount("/api/users", router)

	return mainRouter
}

func shutdown() {
	testDB := mongoTestClient.Database(mongoTestDBName)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	filter := bson.M{}
	testDB.Collection("users").DeleteMany(ctx, filter)
}
