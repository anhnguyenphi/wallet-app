package share

import (
	"context"
	"database/sql"
)

type (
	DB interface {
		BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
		SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	}

	TX interface {
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		Rollback() error
		Commit() error
	}
)
