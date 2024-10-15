package repository

import (
	"context"
	"database/sql"

	"github.com/replu/sqlc-read-write-sample/internal/model"
	"github.com/replu/sqlc-read-write-sample/internal/repository/sqlc"
	"github.com/replu/sqlc-read-write-sample/internal/util/database"
)

type Repository struct {
	dba     *database.Accessor
	queries *sqlc.Queries
}

func NewRepository(dba *database.Accessor, conn *sql.DB) *Repository {
	return &Repository{
		dba:     dba,
		queries: sqlc.New(conn),
	}
}

func (r *Repository) Get(ctx context.Context, id int64) (*model.User, error) {
	tx, err := r.dba.GetTxFromContext(ctx)
	if err != nil {
		return nil, err
	}
	q := r.queries
	if tx != nil {
		q = q.WithTx(tx)
	}

	u, err := q.UserGet(ctx, uint64(id))
	if err != nil {
		return nil, err
	}
	return &model.User{ID: int64(u.ID), Name: u.Name}, nil
}

func (r *Repository) Create(ctx context.Context, name string) (*model.User, error) {
	tx, err := r.dba.GetTxFromContext(ctx)
	if err != nil {
		return nil, err
	}
	q := r.queries
	if tx != nil {
		q = q.WithTx(tx)
	}

	res, err := r.queries.UserCreate(ctx, name)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	u, err := r.queries.UserGet(ctx, uint64(id))
	if err != nil {
		return nil, err
	}

	return &model.User{ID: int64(u.ID), Name: u.Name}, nil
}
