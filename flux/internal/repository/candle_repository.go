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
		return nil, fmt.Errorf("failed to retrieve last candle: %w", err)
	}

	candle.Timestamp = candle.Timestamp.UTC()

	return &candle, nil
}
