package service

import (
	"context"

	"github.com/replu/sqlc-read-write-sample/internal/model"
	"github.com/replu/sqlc-read-write-sample/internal/repository"
	"github.com/replu/sqlc-read-write-sample/internal/util/database"
)

type Service struct {
	dba  *database.Accessor
	repo *repository.Repository
}

func NewService(dba *database.Accessor, repo *repository.Repository) *Service {
	return &Service{
		dba:  dba,
		repo: repo,
	}
}

func (s *Service) Get(ctx context.Context, id int64) (*model.User, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) GetWithTx(ctx context.Context, id int64) (*model.User, error) {
	var user *model.User
	err := s.dba.Transaction(ctx, func(txCtx context.Context) error {
		var err error
		user, err = s.repo.Get(txCtx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Create(ctx context.Context, name string) (*model.User, error) {
	return s.repo.Create(ctx, name)
}
