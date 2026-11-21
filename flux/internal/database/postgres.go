package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"flux/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	poolConfig, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, errors.New("failed to initialize PostgreSQL configuration")
	}

	poolConfig.ConnConfig.Host = cfg.Host
	poolConfig.ConnConfig.Port = uint16(cfg.Port)
	poolConfig.ConnConfig.Database = cfg.Name
	poolConfig.ConnConfig.User = cfg.User
	poolConfig.ConnConfig.Password = cfg.Password

	switch cfg.SSLMode {
	case "disable":
		poolConfig.ConnConfig.TLSConfig = nil
		poolConfig.ConnConfig.Fallbacks = nil
	default:
		return nil, fmt.Errorf("PostgreSQL SSL mode %q is not yet supported", cfg.SSLMode)
	}

	poolConfig.MaxConns = 5

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, errors.New("failed to initialize PostgreSQL connection pool")
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, errors.New("unable to connect to PostgreSQL")
	}

	return pool, nil
}
