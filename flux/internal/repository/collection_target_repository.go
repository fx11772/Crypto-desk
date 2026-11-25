package repository

import (
	"context"
	"fmt"

	"flux/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

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

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, translatePersistenceError(
			fmt.Errorf("failed to retrieve enabled targets: %w", err),
		)
	}

	defer rows.Close()

	targets := make([]domain.CollectionTarget, 0)

	for rows.Next() {
		var target domain.CollectionTarget

		err := rows.Scan(
			&target.ID,
			&target.Exchange,
			&target.Symbol,
			&target.Timeframe,
			&target.Enabled,
			&target.CreatedAt,
			&target.UpdatedAt,
		)

		if err != nil {
			return nil, translatePersistenceError(
				fmt.Errorf("failed to scan collection target: %w", err),
			)
		}

		target.CreatedAt = target.CreatedAt.UTC()
		target.UpdatedAt = target.UpdatedAt.UTC()

		targets = append(targets, target)
	}

	if err := rows.Err(); err != nil {
		return nil, translatePersistenceError(
			fmt.Errorf("failed to iterate collection targets: %w", err),
		)
	}

	return targets, nil
}
