package main

import (
	"context"
	"database/sql"
	"dtdao/greenlight/internal/data"
	"dtdao/greenlight/internal/jsonlog"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string

	db struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type application struct {
	config config
	models data.Models
	logger *jsonlog.Logger
}

func main() {
	var cfg config

	// logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	flag.IntVar(&cfg.port, "port", 4000, "Api server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", fmt.Sprintf("%s%s", os.Getenv("GREENLIGHT_DB_DSN"), "?sslmode=disable"), "PostgrSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSql max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSql max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSql max connections idle time")

	flag.Parse()
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)

	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	router := app.routes()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      router,
		ErrorLog:     log.New(logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.PrintInfo("Starting server", map[string]string{"addr": srv.Addr, "env": cfg.env})

	err = srv.ListenAndServe()
	logger.PrintFatal(err, nil)
}

func openDB(cfg config) (*sql.DB, error) {

	db, err := sql.Open("postgres", cfg.db.dsn)
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
