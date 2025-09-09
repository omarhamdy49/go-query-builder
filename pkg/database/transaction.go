package database

import (
	"context"
	"fmt"

	"github.com/omarhamdy49/go-query-builder/pkg/types"
	"github.com/jmoiron/sqlx"
)

// Transaction wraps a database transaction with additional functionality.
type Transaction struct {
	tx     *sqlx.Tx
	driver types.Driver
}

// NewTransaction creates a new Transaction wrapper.
func NewTransaction(tx *sqlx.Tx, driver types.Driver) *Transaction {
	return &Transaction{
		tx:     tx,
		driver: driver,
	}
}

// QueryContext executes a query that returns rows within the transaction context.
func (t *Transaction) QueryContext(ctx context.Context, query string, args ...interface{}) (types.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

// QueryRowContext executes a query that is expected to return at most one row within the transaction.
func (t *Transaction) QueryRowContext(ctx context.Context, query string, args ...interface{}) types.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}

// ExecContext executes a query without returning any rows within the transaction.
func (t *Transaction) ExecContext(ctx context.Context, query string, args ...interface{}) (types.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

// Begin returns an error as nested transactions are not supported.
func (t *Transaction) Begin() (types.Tx, error) {
	return nil, fmt.Errorf("cannot start a transaction within a transaction")
}

// BeginTx returns an error as nested transactions are not supported.
func (t *Transaction) BeginTx(_ context.Context, _ *types.TxOptions) (types.Tx, error) {
	return nil, fmt.Errorf("cannot start a transaction within a transaction")
}

// Commit commits the transaction.
func (t *Transaction) Commit() error {
	return t.tx.Commit()
}

// Rollback aborts the transaction.
func (t *Transaction) Rollback() error {
	return t.tx.Rollback()
}

// Driver returns the database driver type for this transaction.
func (t *Transaction) Driver() types.Driver {
	return t.driver
}