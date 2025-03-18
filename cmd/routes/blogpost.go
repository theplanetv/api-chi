package routes

import (
	"api-chi/cmd/controllers"
	"api-chi/cmd/middlewares"

	"github.com/go-chi/chi/v5"
)

func BlogPostRoutes(r chi.Router) {
	controller := controllers.BlogPostController{}
	authMiddleware := middlewares.AuthMiddleware{}

	r.Route("/blog/posts", func(r chi.Router) {
		r.Get("/count", controller.Count)
		r.Get("/", controller.GetAll)

		r.With(authMiddleware.CheckLogin).Post("/", controller.Create)
		r.With(authMiddleware.CheckLogin).Patch("/", controller.Update)
		r.With(authMiddleware.CheckLogin).Delete("/{id}", controller.Remove)
	})
}
