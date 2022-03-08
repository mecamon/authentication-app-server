package auth

import (
	"github.com/go-chi/chi/v5"
)

func Routes() *chi.Mux {

	router := chi.NewRouter()

	router.Post("/login", handlers.login)
	router.Post("/register", handlers.register)
	router.Get("/github-access", handlers.accessRequestForGithub)
	router.Post("/github-login", handlers.loginWithGithub)

	return router
}
