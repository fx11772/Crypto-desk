package database

import (
	"context"
	"strings"
	"testing"

	"flux/internal/config"
)

func TestConnectInvalidConfiguration(t *testing.T) {
	cfg := config.DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		Name:     "flux",
		User:     "flux",
		Password: "secret-test-password",
		SSLMode:  "invalid",
	}

	pool, err := Connect(context.Background(), cfg)

	if pool != nil {
		pool.Close()
		t.Fatal("Expected no connection pool")
	}

	if err == nil {
		t.Fatal("Expected a configuration error")
	}

	if !strings.Contains(err.Error(), "FLUX_DB_SSLMODE") {
		t.Errorf("Expected SSL configuration error, got: %v", err)
	}

	if strings.Contains(err.Error(), cfg.Password) {
		t.Fatal("Password exposed in error message")
	}
}

func TestConnectFailedConnectionDoesNotExposePassword(t *testing.T) {
	cfg := config.DatabaseConfig{
		Host:     "127.0.0.1",
		Port:     1,
		Name:     "flux",
		User:     "flux",
		Password: "secret-test-password",
		SSLMode:  "disable",
	}

	pool, err := Connect(context.Background(), cfg)

	if pool != nil {
		pool.Close()
		t.Fatal("Expected no connection pool")
	}

	if err == nil {
		t.Fatal("Expected a connection error")
	}

	if strings.Contains(err.Error(), cfg.Password) {
		t.Fatal("Password exposed in connection error")
	}

	if err.Error() != "unable to connect to PostgreSQL" {
		t.Errorf("Unexpected connection error: %v", err)
	}
}
