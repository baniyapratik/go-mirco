package main

import (
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := APIResponse{
		Error:   false,
		Message: "Hi from broker",
	}
	app.writeJSON(w, http.StatusOK, payload)
}
