package services

import (
	"accountant/internal/domain"
	"accountant/internal/infra/uow"
	"context"
	"fmt"
)

type AccountService struct {
	Uow *uow.UnitOfWork
}

func (s *AccountService) List(ctx context.Context, userId int) ([]domain.Account, error) {
	l, err := uow.Execute(ctx, s.Uow, func(c *uow.Context) ([]domain.Account, error) {
		return c.AccountRepository.ListByUserId(ctx, userId)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}

	return l, nil
}
