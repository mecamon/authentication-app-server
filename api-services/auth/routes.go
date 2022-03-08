package auth

import (
	"github.com/go-chi/chi/v5"
)

func Routes() *chi.Mux {

	router := chi.NewRouter()

	//TODO add a custom appHandler so you can handle errors better. hint: check error handling in https://astaxie.gitbooks.io/
	router.Post("/login", handlers.login)
	router.Post("/register", handlers.register)
	router.Get("/github-access", handlers.accessRequestForGithub)
	router.Post("/github-login", handlers.loginWithGithub)

	return router
}
