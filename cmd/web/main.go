package main

import (
	"database/sql"
	"dtdao/greenlight/internal/data"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type application struct {
	Movies data.MovieModel

	config config
}

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "Api server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", fmt.Sprintf("%s%s", os.Getenv("GREENLIGHT_DB_DSN"), "?sslmode=disable"), "PostgrSQL DSN")

	db, err := sql.Open("postgres", cfg.db.dsn)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	app := &application{
		Movies: data.MovieModel{DB: db},
		config: cfg,
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}

	log.Printf("Starting server on: 8080")

	err = srv.ListenAndServe()
	log.Fatal(err)
}
