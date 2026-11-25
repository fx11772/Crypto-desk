package repository

import (
	"context"
	"errors"
	"io"
	"net"

	"flux/internal/domain"

	"github.com/jackc/pgx/v5/pgconn"
)

func translatePersistenceError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, context.Canceled) {
		return &domain.PersistencePermanentError{Err: err}
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return &domain.PersistenceTemporaryError{Err: err}
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "08000", // Connection exception
			"08001", // Unable to establish connection
			"08003", // Connection does not exist
			"08006", // Connection failure
			"40001", // Serialization failure
			"40P01", // Deadlock detected
			"53300", // Too many connections
			"55P03", // Lock not available
			"57P01", // Admin shutdown
			"57P02", // Crash shutdown
			"57P03": // Cannot connect now

			return &domain.PersistenceTemporaryError{Err: err}

		default:
			return &domain.PersistencePermanentError{Err: err}
		}
	}

	if pgconn.Timeout(err) {
		return &domain.PersistenceTemporaryError{Err: err}
	}

	var netErr net.Error

	if errors.As(err, &netErr) {
		return &domain.PersistenceTemporaryError{Err: err}
	}

	if errors.Is(err, io.EOF) ||
		errors.Is(err, io.ErrUnexpectedEOF) {
		return &domain.PersistenceTemporaryError{Err: err}
	}

	return &domain.PersistencePermanentError{Err: err}
}
