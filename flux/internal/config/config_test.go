package config

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestLoadMissingDatabasePassword(t *testing.T) {
	t.Setenv("FLUX_DB_PASSWORD", "")

	_, err := Load()

	if err == nil {
		t.Fatal("Expected an error when database password is empty")
	}

	if !strings.Contains(err.Error(), "FLUX_DB_PASSWORD") {
		t.Errorf("Expected FLUX_DB_PASSWORD error, got: %v", err)
	}
}

func TestLoadInvalidDatabasePort(t *testing.T) {
	t.Setenv("FLUX_DB_PORT", "invalid")

	_, err := Load()

	if err == nil {
		t.Fatal("Expected an error for an invalid database port")
	}

	if !strings.Contains(err.Error(), "FLUX_DB_PORT") {
		t.Errorf("Expected FLUX_DB_PORT error, got: %v", err)
	}
}

func TestLoadDefaultValues(t *testing.T) {
	t.Setenv("FLUX_DB_PASSWORD", "test-password")
	t.Setenv("FLUX_WORKER_COUNT", "3")
	t.Setenv("FLUX_CANDLE_CLOSE_DELAY_SECONDS", "10")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Expected valid configuration, got: %v", err)
	}

	if cfg.Flux.WorkerCount != 3 {
		t.Errorf("Expected 3 workers, got %d", cfg.Flux.WorkerCount)
	}

	if cfg.Flux.CandleCloseDelay != 10*time.Second {
		t.Errorf("Expected 10 seconds delay, got %v", cfg.Flux.CandleCloseDelay)
	}
}

func TestLoadDefaultWorkerCount(t *testing.T) {
	t.Setenv("FLUX_DB_PASSWORD", "test-password")

	// Remove the variable to test its default value.
	original, exists := os.LookupEnv("FLUX_WORKER_COUNT")
	os.Unsetenv("FLUX_WORKER_COUNT")

	t.Cleanup(func() {
		if exists {
			os.Setenv("FLUX_WORKER_COUNT", original)
		} else {
			os.Unsetenv("FLUX_WORKER_COUNT")
		}
	})

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Expected valid configuration, got: %v", err)
	}

	if cfg.Flux.WorkerCount != 3 {
		t.Errorf("Expected default worker count 3, got %d", cfg.Flux.WorkerCount)
	}
}

func TestLoadInvalidKrakenURL(t *testing.T) {
	t.Setenv("FLUX_DB_PASSWORD", "test-password")
	t.Setenv("FLUX_KRAKEN_BASE_URL", "http://api.kraken.com")

	_, err := Load()

	if err == nil {
		t.Fatal("Expected an error for an invalid Kraken URL")
	}

	if !strings.Contains(err.Error(), "FLUX_KRAKEN_BASE_URL") {
		t.Errorf("Expected FLUX_KRAKEN_BASE_URL error, got: %v", err)
	}
}

func TestLoadInvalidHistoryStart(t *testing.T) {
	t.Setenv("FLUX_DB_PASSWORD", "test-password")
	t.Setenv("FLUX_HISTORY_START", "2026-10-01")

	_, err := Load()

	if err == nil {
		t.Fatal("Expected an error for an invalid history start")
	}

	if !strings.Contains(err.Error(), "FLUX_HISTORY_START") {
		t.Errorf("Expected FLUX_HISTORY_START error, got: %v", err)
	}
}

func TestLoadValidHistoryStart(t *testing.T) {
	t.Setenv("FLUX_DB_PASSWORD", "test-password")
	t.Setenv("FLUX_HISTORY_START", "2026-10-01T08:00:00-04:00")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Expected valid configuration, got: %v", err)
	}

	if cfg.Flux.HistoryStart == nil {
		t.Fatal("Expected HistoryStart to be defined")
	}

	expected := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)

	if !cfg.Flux.HistoryStart.Equal(expected) {
		t.Errorf(
			"Expected %v, got %v",
			expected,
			cfg.Flux.HistoryStart,
		)
	}
}
