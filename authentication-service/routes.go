package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"net/http"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(cors.Default().Handler)

	mux.Use(middleware.Heartbeat("/health"))
	mux.Post("/authenticate", app.Authenticate)
	return mux
}
