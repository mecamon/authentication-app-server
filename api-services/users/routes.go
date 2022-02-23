package users

import (
	"github.com/authentication-app-server/middlewares"
	"github.com/go-chi/chi/v5"
)

func Routes() *chi.Mux {

	router := chi.NewRouter()
	router.Use(middlewares.TokenValidation)
	router.Post("/update", handlers.updateUserInfo)

	return router
}
