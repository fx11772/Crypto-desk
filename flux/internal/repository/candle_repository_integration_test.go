//go:build integration

package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

func setupTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	if os.Getenv("FLUX_DB_HOST") != "flux-db-test" ||
		os.Getenv("FLUX_DB_NAME") != "flux_test" {
		t.Fatal("integration tests require the isolated test database")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	dsn := "postgres://flux_test:flux_test@flux-db-test:5432/flux_test?sslmode=disable"

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to create test pool: %v", err)
	}

	t.Cleanup(pool.Close)

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	return pool
}

func TestCandleRepository_GetLast_NoCandle(t *testing.T) {
	// Arrange
	pool := setupTestPool(t)
	repo := NewCandleRepository(pool)

	ctx := context.Background()

	// Act
	candle, err := repo.GetLast(
		ctx,
		"integration-test",
		"NO-CANDLE/USD",
		"15m",
	)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if candle != nil {
		t.Fatal("expected no candle")
	}
}

func clearCandles(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	_, err := pool.Exec(
		context.Background(),
		"DELETE FROM candles",
	)

	if err != nil {
		t.Fatalf("failed to clear candles: %v", err)
	}
}

func TestCandleRepository_GetLast_ReturnsLatestCandle(t *testing.T) {
	// Arrange
	ctx := context.Background()
	pool := setupTestPool(t)
	repo := NewCandleRepository(pool)

	clearCandles(t, pool)

	t.Cleanup(func() {
		clearCandles(t, pool)
	})

	const query = `
		INSERT INTO candles (
			exchange, symbol, timeframe, timestamp,
			open, high, low, close, volume
		)
		VALUES
			('kraken', 'BTC/USD', '15m',
			 '2026-01-01 12:00:00+00', 100, 110, 90, 105, 10),

			('kraken', 'BTC/USD', '15m',
			 '2026-01-01 12:15:00+00', 105, 115, 95, 110, 20),

			('kraken', 'BTC/USD', '15m',
			 '2026-01-01 12:30:00+00',
			 110.123456789123456789, 120, 100, 115, 30),

			('kraken', 'ETH/USD', '15m',
			 '2026-01-01 12:45:00+00', 200, 210, 190, 205, 10),

			('kraken', 'BTC/USD', '1h',
			 '2026-01-01 13:00:00+00', 300, 310, 290, 305, 10),

			('coinbase', 'BTC/USD', '15m',
			 '2026-01-01 13:15:00+00', 400, 410, 390, 405, 10)
	`

	if _, err := pool.Exec(ctx, query); err != nil {
		t.Fatalf("failed to insert test candles: %v", err)
	}

	// Act
	candle, err := repo.GetLast(
		ctx,
		"kraken",
		"BTC/USD",
		"15m",
	)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if candle == nil {
		t.Fatal("expected a candle")
	}

	expectedTimestamp := time.Date(
		2026, time.January, 1,
		12, 30, 0, 0,
		time.UTC,
	)

	if !candle.Timestamp.Equal(expectedTimestamp) {
		t.Errorf(
			"expected timestamp %v, got %v",
			expectedTimestamp,
			candle.Timestamp,
		)
	}

	if candle.Timestamp.Location() != time.UTC {
		t.Error("expected timestamp in UTC")
	}

	if candle.Exchange != "kraken" ||
		candle.Symbol != "BTC/USD" ||
		candle.Timeframe != "15m" {
		t.Error("unexpected collection target")
	}

	expectedOpen := decimal.RequireFromString(
		"110.123456789123456789",
	)

	if !candle.Open.Equal(expectedOpen) {
		t.Errorf(
			"expected open %s, got %s",
			expectedOpen,
			candle.Open,
		)
	}

	if !candle.High.Equal(decimal.NewFromInt(120)) {
		t.Errorf("expected high 120, got %s", candle.High)
	}

	if !candle.Low.Equal(decimal.NewFromInt(100)) {
		t.Errorf("expected low 100, got %s", candle.Low)
	}

	if !candle.Close.Equal(decimal.NewFromInt(115)) {
		t.Errorf("expected close 115, got %s", candle.Close)
	}

	if !candle.Volume.Equal(decimal.NewFromInt(30)) {
		t.Errorf("expected volume 30, got %s", candle.Volume)
	}
}

func TestCandleRepository_GetLast_DatabaseError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	pool := setupTestPool(t)
	repo := NewCandleRepository(pool)

	pool.Close()

	// Act
	candle, err := repo.GetLast(
		ctx,
		"kraken",
		"BTC/USD",
		"15m",
	)

	// Assert
	if err == nil {
		t.Fatal("expected a database error")
	}

	if candle != nil {
		t.Error("expected no candle when database access fails")
	}
}