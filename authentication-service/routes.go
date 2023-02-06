package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"net/http"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(cors.Default().Handler)

	return mux
}
