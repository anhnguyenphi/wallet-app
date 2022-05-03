package dbclient

import (
	"context"
	"database/sql"
)

type (
	DB interface {
		BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
		QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	}

	TX interface {
		QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		Rollback() error
		Commit() error
	}
	NoOpTx struct {

	}
)

func (n NoOpTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return nil
}

func (n NoOpTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (n NoOpTx) Rollback() error {
	return nil
}

func (n NoOpTx) Commit() error {
	return nil
}

