package users

import (
	"github.com/authentication-app-server/middlewares"
	"github.com/go-chi/chi/v5"
)

func Routes() *chi.Mux {

	//TODO remove photo service
	router := chi.NewRouter()
	router.Use(middlewares.TokenValidation)
	router.Put("/update", handlers.updateUserInfo)
	router.Get("/info", handlers.userInfo)

	return router
}
