package main

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

/* The reciever is defined in /cmd/api/main.go */
func (app *Config) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	r.Use(middleware.Heartbeat("/ping"))

	/* Test endpoint */
	r.Post("/test", app.BrokerTest)

	/* Single entrypoint for all requests */
	r.Post("/handle", app.HandleRequest)

	return r
}
