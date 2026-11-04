package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"flux/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	// Load local .env if available, without overriding
	// existing environment variables.
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No local .env loaded; using environment variables")
	}

	// Load and validate application configuration.
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Invalid configuration: ", err)
	}

	log.Println("Flux configuration loaded successfully")
	log.Println("Environment:", cfg.Environment)
	log.Println("Database host:", cfg.Database.Host)
	log.Println("Kraken URL:", cfg.Kraken.BaseURL)
	log.Println("Worker count:", cfg.Flux.WorkerCount)

	// Wait for a termination signal.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	log.Println("Flux started")

	<-ctx.Done()

	log.Println("Flux shutting down")
}
