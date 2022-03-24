package main

import (
	"github.com/authentication-app-server/api-services/auth"
	"github.com/authentication-app-server/api-services/users"
	"github.com/authentication-app-server/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

var dbName = "authentication"

func makeRouter(dbConn *db.DB) *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"X-count"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	auth.NewHandlers(dbConn.Client, dbName)
	authRoutes := auth.Routes()

	users.NewHandlers(dbConn.Client, dbName)
	userRoutes := users.Routes()

	router.Mount("/api/auth", authRoutes)
	router.Mount("/api/users", userRoutes)

	return router
}
