package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	cacheTTL         = 10 * time.Minute
	cacheRetryAfter  = time.Minute
	cacheMaxStaleAge = 24 * time.Hour
	shutdownTimeout  = 10 * time.Second
)

func main() {
	username := os.Getenv("GEONAMES_USERNAME")
	if username == "" {
		username = "hsample"
	}

	endpoint := os.Getenv("GEONAMES_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://secure.geonames.org/searchJSON"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	client := &geoNamesClient{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		endpoint:   endpoint,
		username:   username,
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           newHandler(newCityCache(client, cacheTTL, cacheRetryAfter, cacheMaxStaleAge)),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)

	go func() {
		log.Printf("Go API listening on %s", server.Addr)
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		log.Fatal(err)
	case <-ctx.Done():
	}

	stop()
	log.Print("Shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)
	}
}
