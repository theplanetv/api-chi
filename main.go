package main

import (
	"api-chi/cmd/config"
	"api-chi/cmd/routes"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// Load configurations for the API and Web
	config.LoadAuthConfig()
	config.LoadApiConfig()
	config.LoadWebConfig()

	// Create a new router
	r := chi.NewRouter()

	// Use middleware for logging requests
	r.Use(middleware.Logger)

	// CORS settings - restrict only to the allowed origins
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://" + config.WEB_HOST + ":" + config.WEB_PORT},       // Restrict to local dev
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},               // Allow necessary HTTP methods
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}, // Explicitly allow headers
		ExposedHeaders:   []string{"Link"},                                                    // Headers that can be exposed to the frontend
		AllowCredentials: true,                                                                // Allow cookies and other credentials
	}))

	// Define the /api route and its subroutes
	r.Route("/api", func(r chi.Router) {
		routes.AuthRoutes(r)
		routes.BlogPostRoutes(r)
		routes.BlogTagRoutes(r)
	})

	// Start the HTTP server on the specified API port
	err := http.ListenAndServe(":"+config.API_PORT, r)
	if err != nil {
		log.Fatal(err)
	}
}
