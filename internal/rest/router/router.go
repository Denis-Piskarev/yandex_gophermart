package router

import (
	"github.com/DenisquaP/yandex_gophermart/internal"
	"github.com/DenisquaP/yandex_gophermart/internal/rest/endpoints"
	"github.com/DenisquaP/yandex_gophermart/internal/rest/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func NewRouterWithMiddleware(endpoints *endpoints.Endpoints, services *internal.Service) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger) // middleware to logging request

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", endpoints.RegisterUser)
		r.Post("/login", endpoints.LoginUser)

		r.Route("/", func(r chi.Router) {
			r.Use(middlewares.IsAuthorized(services))

			r.Get("/withdrawals", endpoints.GetWithdrawals)

			r.Route("/orders", func(r chi.Router) {
				r.Get("/", endpoints.GetOrders)
				r.Post("/", endpoints.UploadOrder)
			})

			r.Route("/balance", func(r chi.Router) {
				// add orders router
				r.Get("/", endpoints.GetBalance)
			})
		})
	})

	return r
}
