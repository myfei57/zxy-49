package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"venueops/internal/console"
	"venueops/internal/store"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "listen address")
	dataDir := flag.String("data", "./data", "data directory")
	webDir := flag.String("web", "./web", "web template directory")
	seed := flag.Bool("seed", false, "seed demo data when the store is empty")
	flag.Parse()

	storeSvc, err := store.New(*dataDir)
	if err != nil {
		log.Fatalf("store init failed: %v", err)
	}
	server, err := console.NewServer(storeSvc, *webDir)
	if err != nil {
		log.Fatalf("server init failed: %v", err)
	}
	if *seed {
		if err := server.SeedIfEmpty(); err != nil {
			log.Fatalf("seed failed: %v", err)
		}
	}

	httpServer := &http.Server{
		Addr:              *addr,
		Handler:           server.Router(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		log.Printf("venueops listening on %s", *addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
