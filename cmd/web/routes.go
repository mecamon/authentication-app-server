package main

import (
	"github.com/authentication-app-server/api-services/auth"
	"github.com/authentication-app-server/api-services/users"
	"github.com/authentication-app-server/db"
	"github.com/go-chi/chi/v5"
)

var dbName = "authentication"

func makeRouter(dbConn *db.DB) *chi.Mux {
	router := chi.NewRouter()

	auth.NewHandlers(dbConn.Client, dbName)
	authRoutes := auth.Routes()

	users.NewHandlers(dbConn.Client, dbName)
	userRoutes := users.Routes()

	router.Mount("/api/auth", authRoutes)
	router.Mount("/api/users", userRoutes)

	return router
}
