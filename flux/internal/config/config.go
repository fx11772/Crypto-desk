package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config contains all Flux application settings.
type Config struct {
	Environment string
	LogLevel    string
	Database    DatabaseConfig
	Kraken      KrakenConfig
	Flux        FluxConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
}

type KrakenConfig struct {
	BaseURL string
	Timeout time.Duration
}

type FluxConfig struct {
	WorkerCount      int
	CandleCloseDelay time.Duration
	HistoryStart     *time.Time
}

// Load reads and validates the environment configuration.
func Load() (*Config, error) {
	dbPort, err := getInt("FLUX_DB_PORT", 5432)
	if err != nil {
		return nil, err
	}

	timeout, err := getInt("FLUX_KRAKEN_TIMEOUT_SECONDS", 15)
	if err != nil {
		return nil, err
	}

	workers, err := getInt("FLUX_WORKER_COUNT", 3)
	if err != nil {
		return nil, err
	}

	delay, err := getInt("FLUX_CANDLE_CLOSE_DELAY_SECONDS", 10)
	if err != nil {
		return nil, err
	}

	var historyStart *time.Time
	if value := os.Getenv("FLUX_HISTORY_START"); value != "" {
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return nil, fmt.Errorf("invalid FLUX_HISTORY_START: %w", err)
		}
		utc := parsed.UTC()
		historyStart = &utc
	}

	cfg := &Config{
		Environment: getEnv("FLUX_ENV", "development"),
		LogLevel:    getEnv("FLUX_LOG_LEVEL", "info"),
		Database: DatabaseConfig{
			Host:     getEnv("FLUX_DB_HOST", "localhost"),
			Port:     dbPort,
			Name:     getEnv("FLUX_DB_NAME", "flux"),
			User:     getEnv("FLUX_DB_USER", "flux"),
			Password: os.Getenv("FLUX_DB_PASSWORD"),
			SSLMode:  getEnv("FLUX_DB_SSLMODE", "disable"),
		},
		Kraken: KrakenConfig{
			BaseURL: getEnv("FLUX_KRAKEN_BASE_URL", "https://api.kraken.com"),
			Timeout: time.Duration(timeout) * time.Second,
		},
		Flux: FluxConfig{
			WorkerCount:      workers,
			CandleCloseDelay: time.Duration(delay) * time.Second,
			HistoryStart:     historyStart,
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks required settings and valid ranges.
func (c *Config) Validate() error {
	if c.Database.Host == "" ||
		c.Database.Name == "" ||
		c.Database.User == "" {
		return fmt.Errorf("missing PostgreSQL configuration")
	}

	if c.Database.Password == "" {
		return fmt.Errorf("FLUX_DB_PASSWORD is required")
	}

	if c.Database.Port < 1 || c.Database.Port > 65535 {
		return fmt.Errorf("FLUX_DB_PORT must be between 1 and 65535")
	}

	switch c.Database.SSLMode {
	case "disable", "allow", "prefer", "require", "verify-ca", "verify-full":
	default:
		return fmt.Errorf("invalid FLUX_DB_SSLMODE")
	}

	parsedURL, err := url.Parse(c.Kraken.BaseURL)
	if err != nil || parsedURL.Scheme != "https" ||
		parsedURL.Host == "" {
		return fmt.Errorf("FLUX_KRAKEN_BASE_URL must be a valid HTTPS URL")
	}

	if c.Kraken.Timeout <= 0 {
		return fmt.Errorf("Kraken timeout must be positive")
	}

	if c.Flux.WorkerCount <= 0 {
		return fmt.Errorf("FLUX_WORKER_COUNT must be positive")
	}

	if c.Flux.CandleCloseDelay < 0 {
		return fmt.Errorf("candle close delay cannot be negative")
	}

	return nil
}

// getEnv returns the value or a default if unset.
func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return strings.TrimSpace(value)
}

// getInt reads an integer environment variable.
func getInt(key string, fallback int) (int, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback, nil
	}

	number, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", key, err)
	}

	return number, nil
}
