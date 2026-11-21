package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"flux/internal/config"
	"flux/internal/database"

	"github.com/joho/godotenv"
)

func main() {
	// Load local .env without overriding existing variables.
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

	// Handle termination signals.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// Initialize PostgreSQL connection.
	pool, err := database.Connect(ctx, cfg.Database)
	if err != nil {
		log.Fatal("PostgreSQL initialization failed: ", err)
	}
	defer pool.Close()

	log.Println("PostgreSQL connection established")
	log.Println("Flux started")

	<-ctx.Done()

	log.Println("Flux shutting down")
}
