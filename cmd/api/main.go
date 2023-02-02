package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/arnab4477/Parkour_API/internal/data"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

// The config struct
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdletime  string
	}
}

// The application struct
type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func main() {

	// An instance of the config struct
	var cfg config

	// Set flags and their default values
	flag.IntVar(&cfg.port, "port", 7001, "The API port")
	flag.StringVar(&cfg.env, "env", "development", "Enviroment (development | staging | production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connetions")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connetions")
	flag.StringVar(&cfg.db.maxIdletime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	flag.Parse()

	// Logger function for customized logging
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Create a database connection pool
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	// Make sure the connection to the pool closes before the main()
	defer db.Close()
	logger.Printf("Database connecton establisted")

	// An instance of the application struct
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	// An HTTP server
	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 40 * time.Second,
	}

	// Start the server
	logger.Printf("Starting %s server on port %s", cfg.env, srv.Addr)

	err = srv.ListenAndServe()
	logger.Fatal(err)
}

// This function returns a SQLdb connection pool
func openDB(cfg config) (*sql.DB, error) {
	// Create an empty connection pool with the config DSN
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// SEt the configuration settings for the connections from the flags
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	// Create a Go time object from the vlue passed as maxIdleTime
	// Set the value to the appropriate settings
	duration, err := time.ParseDuration(cfg.db.maxIdletime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	// Create a context with 5 second time out
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// reate a new connection with the context reatd above
	// If the connecto is not established within 5 seconds, this will return an error
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	// Return the connection pool
	return db, nil
}
