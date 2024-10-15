package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/replu/sqlc-read-write-sample/internal/repository/sqlc"
)

type (
	Accessor struct {
		writerDB *sql.DB
		readerDB *sql.DB
	}
)

func NewAccessor(
	writerDB *sql.DB,
	readerDB *sql.DB,
) *Accessor {
	return &Accessor{
		writerDB: writerDB,
		readerDB: readerDB,
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

func (dba *Accessor) ExecContext(
	ctx context.Context,
	query string,
	args ...interface{},
) (sql.Result, error) {
	tx, err := dba.getTxFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var da sqlc.DBTX
	if tx != nil {
		da = tx
	} else {
		da = dba.writerDB
	}

	result, err := da.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (dba *Accessor) PrepareContext(
	ctx context.Context,
	query string,
) (*sql.Stmt, error) {
	panic("not implemented")
}

func (dba *Accessor) QueryContext(
	ctx context.Context,
	query string,
	args ...interface{},
) (*sql.Rows, error) {
	slog.Info("called QueryContext", slog.String("query", query), slog.Any("args", args))
	tx, err := dba.getTxFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var da sqlc.DBTX
	if tx != nil {
		slog.Info("use writerDB")
		da = tx
	} else {
		slog.Info("use readerDB")
		da = dba.readerDB
	}

	slog.Info("QueryContext", slog.String("query", query), slog.Any("args", args))
	result, err := da.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (dba *Accessor) QueryRowContext(
	ctx context.Context,
	query string,
	args ...interface{},
) *sql.Row {
	tx, err := dba.getTxFromContext(ctx)
	if err != nil {
		// TODO sql.Rowにerrorを設定する方法がわからないため後回し
		return &sql.Row{}
	}

	var da sqlc.DBTX
	if tx != nil {
		slog.Info("use writerDB")
		da = tx
	} else {
		slog.Info("use readerDB")
		da = dba.readerDB
	}
	result := da.QueryRowContext(ctx, query, args...)

	return result
}

type ctxKeyTx struct{}

func (dba *Accessor) withTxContext(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, &ctxKeyTx{}, tx)
}

func (dba *Accessor) getTxFromContext(ctx context.Context) (*sql.Tx, error) {
	if v := ctx.Value(&ctxKeyTx{}); v != nil {
		tx, ok := v.(*sql.Tx)
		if !ok {
			return nil, errors.New("failed to assert *sql.Tx")
		}

		return tx, nil
	}

	return nil, nil
}
