package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"aynshteyn.dev/backend/internal/handler"
	"aynshteyn.dev/backend/internal/middleware"
	"aynshteyn.dev/backend/internal/store"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	cors struct {
		trustedOrigins []string
	}
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "file:subscribers.db", "SQLite data source name")
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = append(cfg.cors.trustedOrigins, val)
		return nil
	})
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := store.NewSQLiteStore(cfg.db.dsn)
	if err != nil {
		logger.Fatalf("cannot open db: %s", err)
	}
	defer db.Close()

	err = db.Migrate()
	if err != nil {
		logger.Fatalf("cannot migrate db: %s", err)
	}

	app := &handler.Application{
		Config: cfg.env,
		Logger: logger,
		Store:  db,
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      middleware.SetupMiddleware(app.Routes(), cfg.cors.trustedOrigins, logger),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on port %d", cfg.env, cfg.port)

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatalf("could not listen on port %d: %s", cfg.port, err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("could not shutdown server: %s", err)
	}
} 