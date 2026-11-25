//go:build integration

package repository

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"flux/internal/domain"

	"github.com/jackc/pgx/v5/pgconn"
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

func TestCandleRepository_InsertBatch_InsertsCandles(t *testing.T) {
	// Arrange
	ctx := context.Background()
	pool := setupTestPool(t)
	repo := NewCandleRepository(pool)

	clearCandles(t, pool)

	t.Cleanup(func() {
		clearCandles(t, pool)
	})

	candles := []domain.Candle{
		{
			Exchange:  "kraken",
			Symbol:    "BTC/USD",
			Timeframe: "15m",
			Timestamp: time.Date(
				2026, time.January, 1, 12, 0, 0, 0, time.UTC,
			),
			Open:   decimal.RequireFromString("100.123456789123456789"),
			High:   decimal.NewFromInt(110),
			Low:    decimal.NewFromInt(90),
			Close:  decimal.NewFromInt(105),
			Volume: decimal.NewFromInt(10),
		},
		{
			Exchange:  "kraken",
			Symbol:    "BTC/USD",
			Timeframe: "15m",
			Timestamp: time.Date(
				2026, time.January, 1, 12, 15, 0, 0, time.UTC,
			),
			Open:   decimal.NewFromInt(105),
			High:   decimal.NewFromInt(115),
			Low:    decimal.NewFromInt(95),
			Close:  decimal.NewFromInt(110),
			Volume: decimal.NewFromInt(20),
		},
	}

	// Act
	result, err := repo.InsertBatch(ctx, candles)

	// Assert
	if err != nil {
		t.Fatalf("unexpected insertion error: %v", err)
	}

	if result.InsertedCount != len(candles) {
		t.Errorf("expected %d inserted candles, got %d", len(candles), result.InsertedCount)
	}

	if result.DuplicateCount != 0 {
		t.Errorf("expected %d duplicates, got %d", 0, result.DuplicateCount)
	}

	var count int

	err = pool.QueryRow(
		ctx,
		"SELECT COUNT(*) FROM candles",
	).Scan(&count)

	if err != nil {
		t.Fatalf("failed to count candles: %v", err)
	}

	if count != len(candles) {
		t.Errorf(
			"expected %d candles, got %d",
			len(candles),
			count,
		)
	}
}

func TestCandleRepository_InsertBatch_EmptyBatch(t *testing.T) {
	// Arrange
	ctx := context.Background()
	pool := setupTestPool(t)
	repo := NewCandleRepository(pool)

	clearCandles(t, pool)

	candles := []domain.Candle{}

	// Act
	result, err := repo.InsertBatch(ctx, candles)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.InsertedCount != 0 {
		t.Errorf("expected %d inserted candles, got %d", 0, result.InsertedCount)
	}

	if result.DuplicateCount != 0 {
		t.Errorf("expected %d duplicates, got %d", 0, result.DuplicateCount)
	}

	var count int

	err = pool.QueryRow(
		ctx,
		"SELECT COUNT(*) FROM candles",
	).Scan(&count)

	if err != nil {
		t.Fatalf("failed to count candles: %v", err)
	}

	if count != 0 {
		t.Errorf("expected 0 candles, got %d", count)
	}
}

func TestCandleRepository_InsertBatch_EmptyBatchWithoutConnection(t *testing.T) {
	// Arrange
	ctx := context.Background()
	pool := setupTestPool(t)
	repo := NewCandleRepository(pool)

	pool.Close()

	// Act
	result, err := repo.InsertBatch(ctx, nil)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.InsertedCount != 0 {
		t.Errorf("expected %d inserted candles, got %d", 0, result.InsertedCount)
	}

	if result.DuplicateCount != 0 {
		t.Errorf("expected %d duplicates, got %d", 0, result.DuplicateCount)
	}
}

func TestCandleRepository_InsertBatch_IgnoresDuplicates(t *testing.T) {
	// Arrange
	ctx := context.Background()
	pool := setupTestPool(t)
	repo := NewCandleRepository(pool)

	clearCandles(t, pool)

	t.Cleanup(func() {
		clearCandles(t, pool)
	})

	original := domain.Candle{
		Exchange:  "kraken",
		Symbol:    "BTC/USD",
		Timeframe: "15m",
		Timestamp: time.Date(
			2026, time.January, 1,
			12, 0, 0, 0, time.UTC,
		),
		Open:   decimal.NewFromInt(100),
		High:   decimal.NewFromInt(110),
		Low:    decimal.NewFromInt(90),
		Close:  decimal.NewFromInt(105),
		Volume: decimal.NewFromInt(10),
	}

	if _, err := repo.InsertBatch(ctx, []domain.Candle{original}); err != nil {
		t.Fatalf("failed to insert original candle: %v", err)
	}

	duplicate := original
	duplicate.Open = decimal.NewFromInt(200)
	duplicate.High = decimal.NewFromInt(210)
	duplicate.Low = decimal.NewFromInt(190)
	duplicate.Close = decimal.NewFromInt(205)
	duplicate.Volume = decimal.NewFromInt(50)

	newCandle := original
	newCandle.Timestamp = original.Timestamp.Add(15 * time.Minute)

	// Act
	result, err := repo.InsertBatch(ctx, []domain.Candle{
		duplicate,
		newCandle,
	})

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.InsertedCount != 1 {
		t.Errorf("expected %d inserted candles, got %d", 1, result.InsertedCount)
	}

	if result.DuplicateCount != 1 {
		t.Errorf("expected %d duplicates, got %d", 1, result.DuplicateCount)
	}

	var count int

	err = pool.QueryRow(
		ctx,
		"SELECT COUNT(*) FROM candles",
	).Scan(&count)

	if err != nil {
		t.Fatalf("failed to count candles: %v", err)
	}

	if count != 2 {
		t.Errorf("expected 2 candles, got %d", count)
	}

	var open, high, low, close, volume decimal.Decimal

	err = pool.QueryRow(
		ctx,
		`SELECT open, high, low, close, volume
		 FROM candles
		 WHERE exchange = $1
		   AND symbol = $2
		   AND timeframe = $3
		   AND timestamp = $4`,
		original.Exchange,
		original.Symbol,
		original.Timeframe,
		original.Timestamp,
	).Scan(&open, &high, &low, &close, &volume)

	if err != nil {
		t.Fatalf("failed to retrieve original candle: %v", err)
	}

	if !open.Equal(original.Open) ||
		!high.Equal(original.High) ||
		!low.Equal(original.Low) ||
		!close.Equal(original.Close) ||
		!volume.Equal(original.Volume) {
		t.Error("original candle was modified")
	}
}

func TestCandleRepository_InsertBatch_PreservesDecimalPrecision(t *testing.T) {
	// Arrange
	ctx := context.Background()
	pool := setupTestPool(t)
	repo := NewCandleRepository(pool)

	clearCandles(t, pool)

	t.Cleanup(func() {
		clearCandles(t, pool)
	})

	candle := domain.Candle{
		Exchange:  "kraken",
		Symbol:    "BTC/USD",
		Timeframe: "15m",
		Timestamp: time.Date(
			2026, time.January, 1,
			12, 0, 0, 0, time.UTC,
		),
		Open:   decimal.RequireFromString("100.123456789123456789"),
		High:   decimal.RequireFromString("110.987654321987654321"),
		Low:    decimal.RequireFromString("90.000000000000000001"),
		Close:  decimal.RequireFromString("105.111111111111111111"),
		Volume: decimal.RequireFromString("0.000000000123456789"),
	}

	// Act
	_, err := repo.InsertBatch(ctx, []domain.Candle{candle})

	// Assert
	if err != nil {
		t.Fatalf("unexpected insertion error: %v", err)
	}

	var open, high, low, close, volume string

	err = pool.QueryRow(
		ctx,
		`SELECT
			open::text,
			high::text,
			low::text,
			close::text,
			volume::text
		FROM candles
		WHERE exchange = $1
		  AND symbol = $2
		  AND timeframe = $3
		  AND timestamp = $4`,
		candle.Exchange,
		candle.Symbol,
		candle.Timeframe,
		candle.Timestamp,
	).Scan(&open, &high, &low, &close, &volume)

	if err != nil {
		t.Fatalf("failed to retrieve inserted candle: %v", err)
	}

	if open != candle.Open.String() {
		t.Errorf("expected open %s, got %s", candle.Open, open)
	}

	if high != candle.High.String() {
		t.Errorf("expected high %s, got %s", candle.High, high)
	}

	if low != candle.Low.String() {
		t.Errorf("expected low %s, got %s", candle.Low, low)
	}

	if close != candle.Close.String() {
		t.Errorf("expected close %s, got %s", candle.Close, close)
	}

	if volume != candle.Volume.String() {
		t.Errorf("expected volume %s, got %s", candle.Volume, volume)
	}
}

func TestCandleRepository_InsertBatch_RollsBackOnError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	pool := setupTestPool(t)
	repo := NewCandleRepository(pool)

	clearCandles(t, pool)

	t.Cleanup(func() {
		clearCandles(t, pool)
	})

	first := domain.Candle{
		Exchange:  "kraken",
		Symbol:    "BTC/USD",
		Timeframe: "15m",
		Timestamp: time.Date(
			2026, time.January, 1,
			12, 0, 0, 0, time.UTC,
		),
		Open:   decimal.NewFromInt(100),
		High:   decimal.NewFromInt(110),
		Low:    decimal.NewFromInt(90),
		Close:  decimal.NewFromInt(105),
		Volume: decimal.NewFromInt(10),
	}

	invalid := first
	invalid.Timestamp = first.Timestamp.Add(15 * time.Minute)
	invalid.Exchange = strings.Repeat("a", 51)

	last := first
	last.Timestamp = first.Timestamp.Add(30 * time.Minute)

	// Act
	result, err := repo.InsertBatch(ctx, []domain.Candle{
		first,
		invalid,
		last,
	})

	// Assert
	if err == nil {
		t.Fatal("expected an insertion error")
	}

	if result.InsertedCount != 0 {
		t.Errorf("expected %d inserted candles, got %d", 0, result.InsertedCount)
	}

	if result.DuplicateCount != 0 {
		t.Errorf("expected %d duplicates, got %d", 0, result.DuplicateCount)
	}

	var pgErr *pgconn.PgError

	if !errors.As(err, &pgErr) || pgErr.Code != "22001" {
		t.Fatalf("expected a string length error, got %v", err)
	}

	var count int

	err = pool.QueryRow(
		ctx,
		"SELECT COUNT(*) FROM candles",
	).Scan(&count)

	if err != nil {
		t.Fatalf("failed to count candles: %v", err)
	}

	if count != 0 {
		t.Errorf(
			"expected 0 candles after rollback, got %d",
			count,
		)
	}
}

func TestCandleRepository_InsertBatch_AllDuplicates(t *testing.T) {
	// Arrange
	ctx := context.Background()
	pool := setupTestPool(t)
	repo := NewCandleRepository(pool)

	clearCandles(t, pool)

	t.Cleanup(func() {
		clearCandles(t, pool)
	})

	first := domain.Candle{
		Exchange:  "kraken",
		Symbol:    "BTC/USD",
		Timeframe: "15m",
		Timestamp: time.Date(
			2026, time.January, 1,
			12, 0, 0, 0, time.UTC,
		),
		Open:   decimal.NewFromInt(100),
		High:   decimal.NewFromInt(110),
		Low:    decimal.NewFromInt(90),
		Close:  decimal.NewFromInt(105),
		Volume: decimal.NewFromInt(10),
	}

	second := first
	second.Timestamp = first.Timestamp.Add(15 * time.Minute)

	candles := []domain.Candle{first, second}

	if _, err := repo.InsertBatch(ctx, candles); err != nil {
		t.Fatalf("failed to insert initial candles: %v", err)
	}

	// Act
	result, err := repo.InsertBatch(ctx, candles)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.InsertedCount != 0 {
		t.Errorf(
			"expected 0 inserted candles, got %d",
			result.InsertedCount,
		)
	}

	if result.DuplicateCount != len(candles) {
		t.Errorf(
			"expected %d duplicates, got %d",
			len(candles),
			result.DuplicateCount,
		)
	}

	var count int

	err = pool.QueryRow(
		ctx,
		"SELECT COUNT(*) FROM candles",
	).Scan(&count)

	if err != nil {
		t.Fatalf("failed to count candles: %v", err)
	}

	if count != len(candles) {
		t.Errorf("expected %d candles, got %d", len(candles), count)
	}
}
