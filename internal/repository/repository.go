package repository

import (
	"context"

	"github.com/replu/sqlc-read-write-sample/internal/model"
	"github.com/replu/sqlc-read-write-sample/internal/repository/sqlc"
	"github.com/replu/sqlc-read-write-sample/internal/util/database"
)

type Repository struct {
	dba     *database.Accessor
	queries *sqlc.Queries
}

func NewRepository(dba *database.Accessor) *Repository {
	return &Repository{
		dba:     dba,
		queries:  sqlc.New(),
	}
}

func (r *Repository) Get(ctx context.Context, id int64) (*model.User, error) {
	u, err := r.queries.UserGet(ctx, r.dba, uint64(id))
	if err != nil {
		return nil, err
	}
	return &model.User{ID: int64(u.ID), Name: u.Name}, nil
}

func (r *Repository) Create(ctx context.Context, name string) (*model.User, error) {
	res, err := r.queries.UserCreate(ctx, r.dba, name)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	u, err := r.queries.UserGet(ctx, r.dba, uint64(id))
	if err != nil {
		return nil, err
	}

	return &model.User{ID: int64(u.ID), Name: u.Name}, nil
}
