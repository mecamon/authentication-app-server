package auth

import (
	"github.com/go-chi/chi/v5"
	"testing"
)

func TestRoutes(t *testing.T) {
	var i interface{}

	i = Routes()

	if _, ok := i.(*chi.Mux); !ok {
		t.Errorf("Wrong type received. Expected *chi.Mux but got another")
	}
}
