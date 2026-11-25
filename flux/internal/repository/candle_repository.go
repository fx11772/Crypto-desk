package repository

import (
	"context"
	"errors"
	"fmt"

	"flux/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CandleRepository struct {
	pool *pgxpool.Pool
}

func NewCandleRepository(pool *pgxpool.Pool) *CandleRepository {
	return &CandleRepository{
		pool: pool,
	}
}

func (r *CandleRepository) GetLast(
	ctx context.Context,
	exchange string,
	symbol string,
	timeframe string,
) (*domain.Candle, error) {
	const query = `
		SELECT
			exchange,
			symbol,
			timeframe,
			timestamp,
			open,
			high,
			low,
			close,
			volume
		FROM candles
		WHERE exchange = $1
			AND symbol = $2
			AND timeframe = $3
		ORDER BY timestamp DESC
		LIMIT 1
	`

	var candle domain.Candle

	err := r.pool.QueryRow(
		ctx,
		query,
		exchange,
		symbol,
		timeframe,
	).Scan(
		&candle.Exchange,
		&candle.Symbol,
		&candle.Timeframe,
		&candle.Timestamp,
		&candle.Open,
		&candle.High,
		&candle.Low,
		&candle.Close,
		&candle.Volume,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, translatePersistenceError(
			fmt.Errorf("failed to retrieve last candle: %w", err),
		)
	}

	candle.Timestamp = candle.Timestamp.UTC()

	return &candle, nil
}

func (r *CandleRepository) InsertBatch(
	ctx context.Context,
	candles []domain.Candle,
) (domain.SaveCandlesResult, error) {
	var result domain.SaveCandlesResult

	if len(candles) == 0 {
		return result, nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.SaveCandlesResult{},
			translatePersistenceError(
				fmt.Errorf("failed to begin candle transaction: %w", err),
			)
	}

	defer tx.Rollback(context.Background())

	const query = `
		INSERT INTO candles (
			exchange,
			symbol,
			timeframe,
			timestamp,
			open,
			high,
			low,
			close,
			volume
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (exchange, symbol, timeframe, timestamp)
		DO NOTHING
	`

	for _, candle := range candles {
		tag, err := tx.Exec(
			ctx,
			query,
			candle.Exchange,
			candle.Symbol,
			candle.Timeframe,
			candle.Timestamp.UTC(),
			candle.Open,
			candle.High,
			candle.Low,
			candle.Close,
			candle.Volume,
		)

		if err != nil {
			return domain.SaveCandlesResult{},
				translatePersistenceError(
					fmt.Errorf("failed to insert candle: %w", err),
				)
		}

		if tag.RowsAffected() == 1 {
			result.InsertedCount++
		} else {
			result.DuplicateCount++
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.SaveCandlesResult{},
			translatePersistenceError(
				fmt.Errorf("failed to commit candle transaction: %w", err),
			)
	}

	return result, nil
}
