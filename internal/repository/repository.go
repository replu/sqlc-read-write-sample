package repository

import (
	"context"
	"database/sql"

	"github.com/replu/sqlc-read-write-sample/internal/model"
	"github.com/replu/sqlc-read-write-sample/internal/repository/sqlc"
)

type Repository struct {
	queries *sqlc.Queries
}

func NewRepository(writerConn, readerConn *sql.DB) *Repository {
	return &Repository{
		queries: sqlc.New(writerConn, readerConn),
	}
}

func (r *Repository) Get(ctx context.Context, id int64) (*model.User, error) {
	u, err := r.queries.UserGet(ctx, uint64(id))
	if err != nil {
		return nil, err
	}
	return &model.User{ID: int64(u.ID), Name: u.Name}, nil
}

func (r *Repository) Create(ctx context.Context, name string) (*model.User, error) {
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
