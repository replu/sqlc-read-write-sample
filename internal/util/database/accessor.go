package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
)

type (
	Accessor struct {
		writerDB *sql.DB
		readerDB *sql.DB
	}
)

func NewAccessor(
	writerDB *sql.DB,
) *Accessor {
	return &Accessor{
		writerDB: writerDB,
	}
}

func (dba *Accessor) Transaction(ctx context.Context, txFunc func(context.Context) error) (err error) {
	tx, err := dba.writerDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			if err := tx.Rollback(); err != nil {
				slog.ErrorContext(ctx, "failed to rollback", slog.String("error", err.Error()))
			}
			slog.InfoContext(ctx, "completed: tx.Rollback")
			err = fmt.Errorf("paniced on execution txFunc: %v", r)
		}
	}()

	txCtx := dba.withTxContext(ctx, tx)

	if txErr := txFunc(txCtx); txErr != nil {
		if err := tx.Rollback(); err != nil {
			slog.ErrorContext(ctx, "failed to rollback", slog.String("error", err.Error()))
		}
		slog.InfoContext(ctx, "completed: tx.Rollback")

		return txErr
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

type ctxKeyTx struct{}

func (dba *Accessor) withTxContext(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, &ctxKeyTx{}, tx)
}

func GetTxFromContext(ctx context.Context) (*sql.Tx, error) {
	if v := ctx.Value(&ctxKeyTx{}); v != nil {
		tx, ok := v.(*sql.Tx)
		if !ok {
			return nil, errors.New("failed to assert *sql.Tx")
		}

		return tx, nil
	}

	return nil, nil
}
