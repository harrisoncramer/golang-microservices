package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "8080"

type Config struct{}

func main() {

	/* The base receiver for our routes */
	app := Config{}

	log.Printf("Starting broker service on %s", webPort)

	/* Define HTTP Server*/
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	/* Start the server */
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
