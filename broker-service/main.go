package main

import (
	"fmt"
	"log"
	"net/http"
)

const PORT = "8081"

func main() {

	app := Config{}

	log.Printf("starting broker-service on port %s\n", PORT)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: app.routes(),
	}
	// start the server
	err := srv.ListenAndServe()
	if err != nil {

		log.Panic(err)
	}
}
