package main

import (
	"github.com/authentication-app-server/api-services/auth"
	"github.com/authentication-app-server/db"
	"github.com/go-chi/chi/v5"
)

func makeRouter(dbConn *db.DB) *chi.Mux {
	router := chi.NewRouter()

	auth.NewHandlers(dbConn.Client, "authentication")
	authRoutes := auth.Routes()

	router.Mount("/api/auth", authRoutes)

	return router
}
