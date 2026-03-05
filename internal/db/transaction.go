package db

import (
	"context"
	"database/sql"
	"fmt"
)

type ctxKey int

const txKey ctxKey = iota

type TransactionManager struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

func (tm *TransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context, tx Querier) error) error {
	tx, err := tm.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	txQuerier := New(tx)
	ctxWithTx := context.WithValue(ctx, txKey, txQuerier)

	if err := fn(ctxWithTx, txQuerier); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback failed after error: %w (original: %v)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func GetQuerierFromContext(ctx context.Context, fallback Querier) Querier {
	if tx, ok := ctx.Value(txKey).(Querier); ok {
		return tx
	}
	return fallback
}
