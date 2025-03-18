package routes

import (
	"api-chi/cmd/controllers"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(r chi.Router) {
	controller := controllers.AuthController{}

	r.Route("/auth", func(r chi.Router) {
		r.Get("/verify", controller.Verify)
		r.Post("/login", controller.Login)
		r.Delete("/logout", controller.Logout)
	})
}
