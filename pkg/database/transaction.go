package database

import (
	"context"
	"fmt"

	"github.com/go-query-builder/querybuilder/pkg/types"
	"github.com/jmoiron/sqlx"
)

type Transaction struct {
	tx     *sqlx.Tx
	driver types.Driver
}

func NewTransaction(tx *sqlx.Tx, driver types.Driver) *Transaction {
	return &Transaction{
		tx:     tx,
		driver: driver,
	}
}

func (t *Transaction) QueryContext(ctx context.Context, query string, args ...interface{}) (types.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

func (t *Transaction) QueryRowContext(ctx context.Context, query string, args ...interface{}) types.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}

func (t *Transaction) ExecContext(ctx context.Context, query string, args ...interface{}) (types.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

func (t *Transaction) Begin() (types.Tx, error) {
	return nil, fmt.Errorf("cannot start a transaction within a transaction")
}

func (t *Transaction) BeginTx(ctx context.Context, opts *types.TxOptions) (types.Tx, error) {
	return nil, fmt.Errorf("cannot start a transaction within a transaction")
}

func (t *Transaction) Commit() error {
	return t.tx.Commit()
}

func (t *Transaction) Rollback() error {
	return t.tx.Rollback()
}

func (t *Transaction) Driver() types.Driver {
	return t.driver
}