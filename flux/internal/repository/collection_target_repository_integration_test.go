//go:build integration

package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"flux/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

func clearCollectionTargets(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	_, err := pool.Exec(
		context.Background(),
		"TRUNCATE TABLE collection_targets RESTART IDENTITY",
	)

	if err != nil {
		t.Fatalf("failed to clear collection targets: %v", err)
	}
}

func insertCollectionTarget(
	t *testing.T,
	pool *pgxpool.Pool,
	exchange string,
	symbol string,
	timeframe string,
	enabled bool,
) int64 {
	t.Helper()

	const query = `
		INSERT INTO collection_targets (
			exchange,
			symbol,
			timeframe,
			enabled
		)
		VALUES ($1, $2, $3, $4)
		RETURNING collection_target_id
	`

	var id int64

	err := pool.QueryRow(
		context.Background(),
		query,
		exchange,
		symbol,
		timeframe,
		enabled,
	).Scan(&id)

	if err != nil {
		t.Fatalf("failed to insert collection target: %v", err)
	}

	return id
}

func TestCollectionTargetRepository_GetEnabled_EmptyTable(
	t *testing.T,
) {
	// Arrange
	ctx := context.Background()
	pool := setupTestPool(t)
	repo := NewCollectionTargetRepository(pool)

	clearCollectionTargets(t, pool)

	t.Cleanup(func() {
		clearCollectionTargets(t, pool)
	})

	// Act
	targets, err := repo.GetEnabled(ctx)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if targets == nil {
		t.Fatal("expected an empty slice, got nil")
	}

	if len(targets) != 0 {
		t.Errorf("expected 0 targets, got %d", len(targets))
	}
}

func TestCollectionTargetRepository_GetEnabled_AllDisabled(
	t *testing.T,
) {
	// Arrange
	ctx := context.Background()
	pool := setupTestPool(t)
	repo := NewCollectionTargetRepository(pool)

	clearCollectionTargets(t, pool)

	t.Cleanup(func() {
		clearCollectionTargets(t, pool)
	})

	insertCollectionTarget(
		t, pool, "kraken", "BTC/USD", "15m", false,
	)

	insertCollectionTarget(
		t, pool, "kraken", "ETH/USD", "1h", false,
	)

	// Act
	targets, err := repo.GetEnabled(ctx)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if targets == nil {
		t.Fatal("expected an empty slice, got nil")
	}

	if len(targets) != 0 {
		t.Errorf("expected 0 targets, got %d", len(targets))
	}
}

func TestCollectionTargetRepository_GetEnabled_MixedTargets(
	t *testing.T,
) {
	// Arrange
	ctx := context.Background()
	pool := setupTestPool(t)
	repo := NewCollectionTargetRepository(pool)

	clearCollectionTargets(t, pool)

	t.Cleanup(func() {
		clearCollectionTargets(t, pool)
	})

	firstID := insertCollectionTarget(
		t, pool, "kraken", "BTC/USD", "15m", true,
	)

	disabledID := insertCollectionTarget(
		t, pool, "kraken", "ETH/USD", "15m", false,
	)

	secondID := insertCollectionTarget(
		t, pool, "kraken", "SOL/USD", "1h", true,
	)

	// Act
	targets, err := repo.GetEnabled(ctx)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(targets) != 2 {
		t.Fatalf("expected 2 targets, got %d", len(targets))
	}

	targetsByID := make(map[int64]domain.CollectionTarget)

	for _, target := range targets {
		targetsByID[target.ID] = target

		if !target.Enabled {
			t.Errorf("target %d should be enabled", target.ID)
		}
	}

	if _, exists := targetsByID[disabledID]; exists {
		t.Error("disabled target should not be returned")
	}

	first, exists := targetsByID[firstID]

	if !exists {
		t.Fatal("first enabled target was not returned")
	}

	if first.Exchange != "kraken" ||
		first.Symbol != "BTC/USD" ||
		first.Timeframe != "15m" {
		t.Errorf("unexpected first target: %+v", first)
	}

	second, exists := targetsByID[secondID]

	if !exists {
		t.Fatal("second enabled target was not returned")
	}

	if second.Exchange != "kraken" ||
		second.Symbol != "SOL/USD" ||
		second.Timeframe != "1h" {
		t.Errorf("unexpected second target: %+v", second)
	}

	for _, target := range targets {
		if target.CreatedAt.IsZero() {
			t.Errorf("target %d has an empty created_at", target.ID)
		}

		if target.UpdatedAt.IsZero() {
			t.Errorf("target %d has an empty updated_at", target.ID)
		}

		if target.CreatedAt.Location() != time.UTC {
			t.Errorf("target %d created_at is not UTC", target.ID)
		}

		if target.UpdatedAt.Location() != time.UTC {
			t.Errorf("target %d updated_at is not UTC", target.ID)
		}
	}
}

func TestCollectionTargetRepository_GetEnabled_DatabaseError(
	t *testing.T,
) {
	// Arrange
	ctx := context.Background()
	pool := setupTestPool(t)
	repo := NewCollectionTargetRepository(pool)

	pool.Close()

	// Act
	targets, err := repo.GetEnabled(ctx)

	// Assert
	if err == nil {
		t.Fatal("expected a persistence error")
	}

	var permanentError *domain.PersistencePermanentError

	if !errors.As(err, &permanentError) {
		t.Fatalf(
			"expected PersistencePermanentError, got %T",
			err,
		)
	}

	if targets != nil {
		t.Error("expected nil targets on error")
	}
}

func TestCollectionTargetRepository_GetEnabled_DeadlineExceeded(
	t *testing.T,
) {
	// Arrange
	pool := setupTestPool(t)
	repo := NewCollectionTargetRepository(pool)

	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(-time.Second),
	)
	defer cancel()

	// Act
	targets, err := repo.GetEnabled(ctx)

	// Assert
	if err == nil {
		t.Fatal("expected a persistence error")
	}

	var temporaryError *domain.PersistenceTemporaryError

	if !errors.As(err, &temporaryError) {
		t.Fatalf(
			"expected PersistenceTemporaryError, got %T",
			err,
		)
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Error("original deadline error was not preserved")
	}

	if targets != nil {
		t.Error("expected nil targets on error")
	}
}

func TestCollectionTargetRepository_GetAll_EmptyTable(
	t *testing.T,
) {
	// Arrange
	ctx := context.Background()
	pool := setupTestPool(t)
	repo := NewCollectionTargetRepository(pool)

	clearCollectionTargets(t, pool)

	t.Cleanup(func() {
		clearCollectionTargets(t, pool)
	})

	// Act
	targets, err := repo.GetAll(ctx)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if targets == nil {
		t.Fatal("expected an empty slice, got nil")
	}

	if len(targets) != 0 {
		t.Errorf("expected 0 targets, got %d", len(targets))
	}
}

func TestCollectionTargetRepository_GetAll_AllDisabled(
	t *testing.T,
) {
	// Arrange
	ctx := context.Background()
	pool := setupTestPool(t)
	repo := NewCollectionTargetRepository(pool)

	clearCollectionTargets(t, pool)

	t.Cleanup(func() {
		clearCollectionTargets(t, pool)
	})

	insertCollectionTarget(
		t, pool, "kraken", "BTC/USD", "15m", false,
	)

	insertCollectionTarget(
		t, pool, "kraken", "ETH/USD", "1h", false,
	)

	// Act
	targets, err := repo.GetAll(ctx)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(targets) != 2 {
		t.Fatalf("expected 2 targets, got %d", len(targets))
	}

	for _, target := range targets {
		if target.Enabled {
			t.Errorf("target %d should be disabled", target.ID)
		}
	}
}

func TestCollectionTargetRepository_GetAll_MixedTargets(
	t *testing.T,
) {
	// Arrange
	ctx := context.Background()
	pool := setupTestPool(t)
	repo := NewCollectionTargetRepository(pool)

	clearCollectionTargets(t, pool)

	t.Cleanup(func() {
		clearCollectionTargets(t, pool)
	})

	firstID := insertCollectionTarget(
		t, pool, "kraken", "BTC/USD", "15m", true,
	)

	secondID := insertCollectionTarget(
		t, pool, "kraken", "ETH/USD", "1h", false,
	)

	thirdID := insertCollectionTarget(
		t, pool, "kraken", "SOL/USD", "15m", true,
	)

	// Act
	targets, err := repo.GetAll(ctx)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(targets) != 3 {
		t.Fatalf("expected 3 targets, got %d", len(targets))
	}

	expected := map[int64]domain.CollectionTarget{
		firstID: {
			ID:        firstID,
			Exchange:  "kraken",
			Symbol:    "BTC/USD",
			Timeframe: "15m",
			Enabled:   true,
		},
		secondID: {
			ID:        secondID,
			Exchange:  "kraken",
			Symbol:    "ETH/USD",
			Timeframe: "1h",
			Enabled:   false,
		},
		thirdID: {
			ID:        thirdID,
			Exchange:  "kraken",
			Symbol:    "SOL/USD",
			Timeframe: "15m",
			Enabled:   true,
		},
	}

	for _, target := range targets {
		expectedTarget, exists := expected[target.ID]

		if !exists {
			t.Errorf("unexpected or duplicate target %d", target.ID)
			continue
		}

		if target.Exchange != expectedTarget.Exchange ||
			target.Symbol != expectedTarget.Symbol ||
			target.Timeframe != expectedTarget.Timeframe ||
			target.Enabled != expectedTarget.Enabled {
			t.Errorf("unexpected target values: %+v", target)
		}

		if target.CreatedAt.IsZero() ||
			target.UpdatedAt.IsZero() {
			t.Errorf("target %d has empty timestamps", target.ID)
		}

		if target.CreatedAt.Location() != time.UTC ||
			target.UpdatedAt.Location() != time.UTC {
			t.Errorf("target %d timestamps are not UTC", target.ID)
		}

		delete(expected, target.ID)
	}

	if len(expected) != 0 {
		t.Errorf("%d expected targets were not returned", len(expected))
	}
}

func TestCollectionTargetRepository_GetByID_EnabledTarget(
	t *testing.T,
) {
	// Arrange
	ctx := context.Background()
	pool := setupTestPool(t)
	repo := NewCollectionTargetRepository(pool)

	clearCollectionTargets(t, pool)

	t.Cleanup(func() {
		clearCollectionTargets(t, pool)
	})

	id := insertCollectionTarget(
		t, pool, "kraken", "BTC/USD", "15m", true,
	)

	insertCollectionTarget(
		t, pool, "kraken", "ETH/USD", "1h", true,
	)

	// Act
	target, err := repo.GetByID(ctx, id)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if target == nil {
		t.Fatal("expected a collection target, got nil")
	}

	if target.ID != id {
		t.Errorf("expected ID %d, got %d", id, target.ID)
	}

	if target.Exchange != "kraken" ||
		target.Symbol != "BTC/USD" ||
		target.Timeframe != "15m" ||
		!target.Enabled {
		t.Errorf("unexpected target values: %+v", target)
	}

	if target.CreatedAt.IsZero() ||
		target.UpdatedAt.IsZero() {
		t.Error("expected non-zero timestamps")
	}

	if target.CreatedAt.Location() != time.UTC ||
		target.UpdatedAt.Location() != time.UTC {
		t.Error("expected UTC timestamps")
	}
}

func TestCollectionTargetRepository_GetByID_DisabledTarget(
	t *testing.T,
) {
	// Arrange
	ctx := context.Background()
	pool := setupTestPool(t)
	repo := NewCollectionTargetRepository(pool)

	clearCollectionTargets(t, pool)

	t.Cleanup(func() {
		clearCollectionTargets(t, pool)
	})

	id := insertCollectionTarget(
		t, pool, "kraken", "ETH/USD", "1h", false,
	)

	// Act
	target, err := repo.GetByID(ctx, id)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if target == nil {
		t.Fatal("expected a disabled collection target, got nil")
	}

	if target.ID != id {
		t.Errorf("expected ID %d, got %d", id, target.ID)
	}

	if target.Exchange != "kraken" ||
		target.Symbol != "ETH/USD" ||
		target.Timeframe != "1h" ||
		target.Enabled {
		t.Errorf("unexpected target values: %+v", target)
	}
}

func TestCollectionTargetRepository_GetByID_NotFound(
	t *testing.T,
) {
	// Arrange
	ctx := context.Background()
	pool := setupTestPool(t)
	repo := NewCollectionTargetRepository(pool)

	clearCollectionTargets(t, pool)

	t.Cleanup(func() {
		clearCollectionTargets(t, pool)
	})

	existingID := insertCollectionTarget(
		t, pool, "kraken", "BTC/USD", "15m", true,
	)

	missingID := existingID + 1

	// Act
	target, err := repo.GetByID(ctx, missingID)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if target != nil {
		t.Errorf("expected nil target, got %+v", target)
	}
}

func TestCollectionTargetRepository_GetAll_DatabaseError(
	t *testing.T,
) {
	// Arrange
	ctx := context.Background()
	pool := setupTestPool(t)
	repo := NewCollectionTargetRepository(pool)

	pool.Close()

	// Act
	targets, err := repo.GetAll(ctx)

	// Assert
	if err == nil {
		t.Fatal("expected a persistence error")
	}

	var permanentError *domain.PersistencePermanentError

	if !errors.As(err, &permanentError) {
		t.Fatalf(
			"expected PersistencePermanentError, got %T",
			err,
		)
	}

	if targets != nil {
		t.Error("expected nil targets on error")
	}
}

func TestCollectionTargetRepository_GetAll_DeadlineExceeded(
	t *testing.T,
) {
	// Arrange
	pool := setupTestPool(t)
	repo := NewCollectionTargetRepository(pool)

	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(-time.Second),
	)
	defer cancel()

	// Act
	targets, err := repo.GetAll(ctx)

	// Assert
	if err == nil {
		t.Fatal("expected a persistence error")
	}

	var temporaryError *domain.PersistenceTemporaryError

	if !errors.As(err, &temporaryError) {
		t.Fatalf(
			"expected PersistenceTemporaryError, got %T",
			err,
		)
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Error("original deadline error was not preserved")
	}

	if targets != nil {
		t.Error("expected nil targets on error")
	}
}

func TestCollectionTargetRepository_GetByID_DatabaseError(
	t *testing.T,
) {
	// Arrange
	ctx := context.Background()
	pool := setupTestPool(t)
	repo := NewCollectionTargetRepository(pool)

	pool.Close()

	// Act
	target, err := repo.GetByID(ctx, 1)

	// Assert
	if err == nil {
		t.Fatal("expected a persistence error")
	}

	var permanentError *domain.PersistencePermanentError

	if !errors.As(err, &permanentError) {
		t.Fatalf(
			"expected PersistencePermanentError, got %T",
			err,
		)
	}

	if target != nil {
		t.Error("expected nil target on error")
	}
}

func TestCollectionTargetRepository_GetByID_DeadlineExceeded(
	t *testing.T,
) {
	// Arrange
	pool := setupTestPool(t)
	repo := NewCollectionTargetRepository(pool)

	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(-time.Second),
	)
	defer cancel()

	// Act
	target, err := repo.GetByID(ctx, 1)

	// Assert
	if err == nil {
		t.Fatal("expected a persistence error")
	}

	var temporaryError *domain.PersistenceTemporaryError

	if !errors.As(err, &temporaryError) {
		t.Fatalf(
			"expected PersistenceTemporaryError, got %T",
			err,
		)
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Error("original deadline error was not preserved")
	}

	if target != nil {
		t.Error("expected nil target on error")
	}
}
