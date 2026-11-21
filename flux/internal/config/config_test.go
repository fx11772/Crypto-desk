package config

import (
	"os"
	"strings"
	"testing"
	"time"
)

func setTestDatabaseEnv(t *testing.T) {
	t.Helper()

	t.Setenv("FLUX_DB_HOST", "localhost")
	t.Setenv("FLUX_DB_PORT", "5432")
	t.Setenv("FLUX_DB_NAME", "flux")
	t.Setenv("FLUX_DB_USER", "flux")
	t.Setenv("FLUX_DB_PASSWORD", "test-password")
	t.Setenv("FLUX_DB_SSLMODE", "disable")
}

func TestLoadMissingDatabasePassword(t *testing.T) {
	setTestDatabaseEnv(t)

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
	setTestDatabaseEnv(t)

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
	setTestDatabaseEnv(t)

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
	setTestDatabaseEnv(t)

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
	setTestDatabaseEnv(t)

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
	setTestDatabaseEnv(t)

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

func TestDatabaseConfigValidation(t *testing.T) {
	tests := []struct {
		name     string
		envKey   string
		value    string
		expected string
	}{
		{"Empty host", "FLUX_DB_HOST", "", "FLUX_DB_HOST"},
		{"Empty database", "FLUX_DB_NAME", "", "FLUX_DB_NAME"},
		{"Empty user", "FLUX_DB_USER", "", "FLUX_DB_USER"},
		{"Empty password", "FLUX_DB_PASSWORD", "", "FLUX_DB_PASSWORD"},
		{"Invalid port", "FLUX_DB_PORT", "abc", "FLUX_DB_PORT"},
		{"Port too low", "FLUX_DB_PORT", "0", "FLUX_DB_PORT"},
		{"Port too high", "FLUX_DB_PORT", "65536", "FLUX_DB_PORT"},
		{"Invalid SSL mode", "FLUX_DB_SSLMODE", "invalid", "FLUX_DB_SSLMODE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("FLUX_DB_HOST", "localhost")
			t.Setenv("FLUX_DB_NAME", "flux")
			t.Setenv("FLUX_DB_USER", "flux")
			t.Setenv("FLUX_DB_PASSWORD", "test-password")
			t.Setenv("FLUX_DB_PORT", "5432")
			t.Setenv("FLUX_DB_SSLMODE", "disable")

			t.Setenv(tt.envKey, tt.value)

			_, err := Load()
			if err == nil {
				t.Fatal("Expected a configuration error")
			}

			if !strings.Contains(err.Error(), tt.expected) {
				t.Errorf("Expected error mentioning %s, got: %v", tt.expected, err)
			}

			if strings.Contains(err.Error(), "test-password") {
				t.Fatal("Database password exposed in error")
			}
		})
	}
}

func TestLoadMissingRequiredDatabaseVariables(t *testing.T) {
	variables := []string{
		"FLUX_DB_HOST",
		"FLUX_DB_NAME",
		"FLUX_DB_USER",
		"FLUX_DB_PASSWORD",
	}

	for _, variable := range variables {
		t.Run(variable, func(t *testing.T) {
			setTestDatabaseEnv(t)

			t.Setenv(variable, "")

			_, err := Load()
			if err == nil {
				t.Fatalf("Expected an error for missing %s", variable)
			}

			if !strings.Contains(err.Error(), variable) {
				t.Errorf(
					"Expected error mentioning %s, got: %v",
					variable,
					err,
				)
			}

			if strings.Contains(err.Error(), "test-password") {
				t.Fatal("Database password exposed in error")
			}
		})
	}
}

func TestLoadUnsetRequiredDatabaseVariables(t *testing.T) {
	variables := []string{
		"FLUX_DB_HOST",
		"FLUX_DB_NAME",
		"FLUX_DB_USER",
		"FLUX_DB_PASSWORD",
	}

	for _, variable := range variables {
		t.Run(variable, func(t *testing.T) {
			setTestDatabaseEnv(t)

			// Save the previous value before removing it.
			original, exists := os.LookupEnv(variable)
			os.Unsetenv(variable)

			t.Cleanup(func() {
				if exists {
					os.Setenv(variable, original)
				} else {
					os.Unsetenv(variable)
				}
			})

			_, err := Load()

			if err == nil {
				t.Fatalf("Expected an error when %s is unset", variable)
			}

			if !strings.Contains(err.Error(), variable) {
				t.Errorf(
					"Expected error mentioning %s, got: %v",
					variable,
					err,
				)
			}

			if strings.Contains(err.Error(), "test-password") {
				t.Fatal("Database password exposed in error")
			}
		})
	}
}
