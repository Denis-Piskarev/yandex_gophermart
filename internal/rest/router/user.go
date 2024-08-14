package router

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (ro *Router) NewUserRouter() http.Handler {
	r := chi.NewRouter()

	r.Route("/user", func(r chi.Router) {
		r.Post("/register", ro.Endpoints.RegisterUser)
	})

	return r
}
