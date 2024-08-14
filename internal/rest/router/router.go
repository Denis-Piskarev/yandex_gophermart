package router

import (
	"github.com/DenisquaP/yandex_gophermart/internal/rest/endpoints"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Router struct {
	Endpoints *endpoints.Endpoints
}

func NewRouter(endpoints *endpoints.Endpoints) *Router {
	return &Router{
		Endpoints: endpoints,
	}
}

func NewRouterWithMiddleware(endpoints *endpoints.Endpoints) http.Handler {
	router := NewRouter(endpoints)

	r := chi.NewRouter()

	r.Mount("/", router.NewUserRouter())

	return r
}
