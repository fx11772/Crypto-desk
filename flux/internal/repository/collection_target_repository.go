package repository

import (
	"context"
	"errors"
	"fmt"

	"flux/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type collectionTargetScanner interface {
	Scan(dest ...any) error
}

type CollectionTargetRepository struct {
	pool *pgxpool.Pool
}

func NewCollectionTargetRepository(
	pool *pgxpool.Pool,
) *CollectionTargetRepository {
	return &CollectionTargetRepository{
		pool: pool,
	}
}

func (r *CollectionTargetRepository) GetEnabled(
	ctx context.Context,
) ([]domain.CollectionTarget, error) {
	const query = `
		SELECT
			collection_target_id,
			exchange,
			symbol,
			timeframe,
			enabled,
			created_at,
			updated_at
		FROM collection_targets
		WHERE enabled = TRUE
	`

	return r.queryTargets(ctx, query)
}

func (r *CollectionTargetRepository) GetAll(
	ctx context.Context,
) ([]domain.CollectionTarget, error) {
	const query = `
		SELECT
			collection_target_id,
			exchange,
			symbol,
			timeframe,
			enabled,
			created_at,
			updated_at
		FROM collection_targets
	`

	return r.queryTargets(ctx, query)
}

func (r *CollectionTargetRepository) GetByID(
	ctx context.Context,
	id int64,
) (*domain.CollectionTarget, error) {
	const query = `
		SELECT
			collection_target_id,
			exchange,
			symbol,
			timeframe,
			enabled,
			created_at,
			updated_at
		FROM collection_targets
		WHERE collection_target_id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)

	target, err := scanCollectionTarget(row)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, translatePersistenceError(
			fmt.Errorf(
				"failed to retrieve collection target %d: %w",
				id,
				err,
			),
		)
	}

	return &target, nil
}

func scanCollectionTarget(
	scanner collectionTargetScanner,
) (domain.CollectionTarget, error) {
	var target domain.CollectionTarget

	err := scanner.Scan(
		&target.ID,
		&target.Exchange,
		&target.Symbol,
		&target.Timeframe,
		&target.Enabled,
		&target.CreatedAt,
		&target.UpdatedAt,
	)

	if err != nil {
		return domain.CollectionTarget{}, err
	}

	target.CreatedAt = target.CreatedAt.UTC()
	target.UpdatedAt = target.UpdatedAt.UTC()

	return target, nil
}

func (r *CollectionTargetRepository) queryTargets(
	ctx context.Context,
	query string,
) ([]domain.CollectionTarget, error) {
	rows, err := r.pool.Query(ctx, query)

	if err != nil {
		return nil, translatePersistenceError(
			fmt.Errorf("failed to retrieve collection targets: %w", err),
		)
	}

	defer rows.Close()

	targets := make([]domain.CollectionTarget, 0)

	for rows.Next() {
		target, err := scanCollectionTarget(rows)

		if err != nil {
			return nil, translatePersistenceError(
				fmt.Errorf("failed to scan collection target: %w", err),
			)
		}

		targets = append(targets, target)
	}

	if err := rows.Err(); err != nil {
		return nil, translatePersistenceError(
			fmt.Errorf("failed to iterate collection targets: %w", err),
		)
	}

	return targets, nil
}
