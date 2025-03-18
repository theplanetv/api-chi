package routes

import (
	"api-chi/cmd/controllers"
	"api-chi/cmd/middlewares"

	"github.com/go-chi/chi/v5"
)

func BlogTagRoutes(r chi.Router) {
	controller := controllers.BlogTagController{}
	authMiddleware := middlewares.AuthMiddleware{}

	r.Route("/blog/tags", func(r chi.Router) {
		r.Get("/count", controller.Count)
		r.Get("/", controller.GetAll)

		r.With(authMiddleware.CheckLogin).Post("/", controller.Create)
		r.With(authMiddleware.CheckLogin).Patch("/", controller.Update)
		r.With(authMiddleware.CheckLogin).Delete("/{id}", controller.Remove)
	})
}
