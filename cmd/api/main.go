package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

// The config struct
type config struct {
	port int
	env string
}


// The application struct
type application struct {
	config config
	logger *log.Logger
}

func main() {

	// An instance of the config struct
	var cfg config

	// Set flags and their default values
	flag.IntVar(&cfg.port, "port", 7000, "The API port")
	flag.StringVar(&cfg.env, "env", "development", "Enviroment (development | staging | production)")
	flag.Parse()

	// Logger function for customized logging
	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)

	// An instance of the application struct
	app := &application {
		config: cfg,
		logger: logger,
	}

	// An HTTP server
	srv := http.Server{
		Addr: fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 40 * time.Second,
	}

	// Start the server
	logger.Printf("Starting %s server on port %s", cfg.env, srv.Addr)

	err := srv.ListenAndServe()
	logger.Fatal(err)
}