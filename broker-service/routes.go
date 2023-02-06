package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"net/http"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()
	// set default cors
	mux.Use(cors.Default().Handler)
	mux.Use(middleware.Heartbeat("/health"))

	mux.Post("/", app.Broker)
	mux.Post("/handle", app.HandleSubmission)

	return mux
}
