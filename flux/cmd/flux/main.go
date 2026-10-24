package main

import (
	"log"

	"github.com/joho/godotenv"
	"flux/internal/config"
)

func main() {
	// Load environment variables from .env
	if err := godotenv.Load("internal/.env"); err != nil {
		log.Fatal("Failed to load .env: ", err)
	}

	// Load and validate application configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Invalid configuration: ", err)
	}

	log.Println("Flux configuration loaded successfully")
	log.Println("Environment:", cfg.Environment)
	log.Println("Database host:", cfg.Database.Host)
	log.Println("Kraken URL:", cfg.Kraken.BaseURL)
	log.Println("Worker count:", cfg.Flux.WorkerCount)
}
