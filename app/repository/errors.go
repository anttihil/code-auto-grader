package repository

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Sentinel errors for common database issues
var (
	// ErrNotFound indicates the requested resource does not exist
	ErrNotFound = errors.New("resource not found")

	// ErrDuplicateKey indicates a unique constraint violation
	ErrDuplicateKey = errors.New("duplicate key violation")

	// ErrForeignKeyViolation indicates a foreign key constraint violation
	ErrForeignKeyViolation = errors.New("foreign key constraint violation")

	// ErrCheckViolation indicates a check constraint violation
	ErrCheckViolation = errors.New("check constraint violation")
)

// PostgreSQL error codes
const (
	pgErrCodeUniqueViolation    = "23505"
	pgErrCodeForeignKeyViolation = "23503"
	pgErrCodeCheckViolation     = "23514"
)

// WrapDBError wraps a database error with appropriate sentinel errors
// This allows handlers to distinguish between 404, 400, and 500 errors
func WrapDBError(err error, resource string) error {
	if err == nil {
		return nil
	}

	// Check for "not found" errors
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w: %s", ErrNotFound, resource)
	}

	// Check for PostgreSQL constraint violations
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgErrCodeUniqueViolation:
			return fmt.Errorf("%w: %s (%s)", ErrDuplicateKey, resource, pgErr.Detail)
		case pgErrCodeForeignKeyViolation:
			return fmt.Errorf("%w: %s (%s)", ErrForeignKeyViolation, resource, pgErr.Detail)
		case pgErrCodeCheckViolation:
			return fmt.Errorf("%w: %s (%s)", ErrCheckViolation, resource, pgErr.Detail)
		}
	}

	// Generic database error
	return fmt.Errorf("database error for %s: %w", resource, err)
}

// IsNotFound checks if an error is a "not found" error
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsDuplicateKey checks if an error is a duplicate key violation
func IsDuplicateKey(err error) bool {
	return errors.Is(err, ErrDuplicateKey)
}

// IsForeignKeyViolation checks if an error is a foreign key violation
func IsForeignKeyViolation(err error) bool {
	return errors.Is(err, ErrForeignKeyViolation)
}

// IsCheckViolation checks if an error is a check constraint violation
func IsCheckViolation(err error) bool {
	return errors.Is(err, ErrCheckViolation)
}
