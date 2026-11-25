package repository

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"testing"

	"flux/internal/domain"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestTranslatePersistenceError_Nil(t *testing.T) {
	if err := translatePersistenceError(nil); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestTranslatePersistenceError_Temporary(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "Connection exception",
			err:  &pgconn.PgError{Code: "08000"},
		},
		{
			name: "Unable to establish connection",
			err:  &pgconn.PgError{Code: "08001"},
		},
		{
			name: "Connection does not exist",
			err:  &pgconn.PgError{Code: "08003"},
		},
		{
			name: "Connection failure",
			err:  &pgconn.PgError{Code: "08006"},
		},
		{
			name: "Serialization failure",
			err:  &pgconn.PgError{Code: "40001"},
		},
		{
			name: "Deadlock",
			err:  &pgconn.PgError{Code: "40P01"},
		},
		{
			name: "Too many connections",
			err:  &pgconn.PgError{Code: "53300"},
		},
		{
			name: "Lock not available",
			err:  &pgconn.PgError{Code: "55P03"},
		},
		{
			name: "Admin shutdown",
			err:  &pgconn.PgError{Code: "57P01"},
		},
		{
			name: "Crash shutdown",
			err:  &pgconn.PgError{Code: "57P02"},
		},
		{
			name: "Server unavailable",
			err:  &pgconn.PgError{Code: "57P03"},
		},
		{
			name: "Deadline exceeded",
			err:  context.DeadlineExceeded,
		},
		{
			name: "Network error",
			err: &net.OpError{
				Op:  "dial",
				Net: "tcp",
				Err: errors.New("connection refused"),
			},
		},
		{
			name: "EOF",
			err:  io.EOF,
		},
		{
			name: "Unexpected EOF",
			err:  io.ErrUnexpectedEOF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := translatePersistenceError(tt.err)

			var temporaryError *domain.PersistenceTemporaryError

			if !errors.As(err, &temporaryError) {
				t.Fatalf(
					"expected PersistenceTemporaryError, got %T",
					err,
				)
			}

			if !errors.Is(err, tt.err) {
				t.Error("original error was not preserved")
			}
		})
	}
}

func TestTranslatePersistenceError_Permanent(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "String too long",
			err:  &pgconn.PgError{Code: "22001"},
		},
		{
			name: "Not null violation",
			err:  &pgconn.PgError{Code: "23502"},
		},
		{
			name: "Unique violation",
			err:  &pgconn.PgError{Code: "23505"},
		},
		{
			name: "Undefined table",
			err:  &pgconn.PgError{Code: "42P01"},
		},
		{
			name: "Syntax error",
			err:  &pgconn.PgError{Code: "42601"},
		},
		{
			name: "Invalid password",
			err:  &pgconn.PgError{Code: "28P01"},
		},
		{
			name: "Context canceled",
			err:  context.Canceled,
		},
		{
			name: "Unknown error",
			err:  errors.New("unknown persistence failure"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := translatePersistenceError(tt.err)

			var permanentError *domain.PersistencePermanentError

			if !errors.As(err, &permanentError) {
				t.Fatalf(
					"expected PersistencePermanentError, got %T",
					err,
				)
			}

			if !errors.Is(err, tt.err) {
				t.Error("original error was not preserved")
			}
		})
	}
}

func TestTranslatePersistenceError_WrappedError(t *testing.T) {
	original := &pgconn.PgError{
		Code: "40P01",
	}

	wrapped := fmt.Errorf(
		"failed to insert candle: %w",
		original,
	)

	err := translatePersistenceError(wrapped)

	var temporaryError *domain.PersistenceTemporaryError

	if !errors.As(err, &temporaryError) {
		t.Fatalf(
			"expected PersistenceTemporaryError, got %T",
			err,
		)
	}

	var pgErr *pgconn.PgError

	if !errors.As(err, &pgErr) {
		t.Fatal("original PostgreSQL error was not preserved")
	}

	if pgErr.Code != "40P01" {
		t.Errorf("unexpected SQLSTATE: %s", pgErr.Code)
	}
}
