package uow

import (
	"accountant/internal/infra/repositories"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UnitOfWork struct {
	pool *pgxpool.Pool
}

type Context struct {
	UserRepository    *repositories.UserRepository
	AccountRepository *repositories.AccountRepository
}

type ExecuteFunc[T any] func(c *Context) (T, error)

func New(pool *pgxpool.Pool) *UnitOfWork {
	return &UnitOfWork{pool: pool}
}

func (u *UnitOfWork) Pool() *pgxpool.Pool {
	return u.pool
}

func Execute[T any](
	ctx context.Context,
	uow *UnitOfWork,
	fn ExecuteFunc[T],
) (T, error) {
	var res T

	c, err := uow.Pool().Acquire(ctx)
	if err != nil {
		return res, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer c.Release()

	res, err = fn(&Context{
		UserRepository: repositories.NewUserRepository(c.Conn()),
	})
	if err != nil {
		return res, fmt.Errorf("failed to execute transaction: %w", err)
	}

	return res, nil
}

func ExecuteTx[T any](
	ctx context.Context,
	uow *UnitOfWork,
	fn ExecuteFunc[T],
	txOptions pgx.TxOptions,
) (T, error) {
	var res T

	c, err := uow.Pool().Acquire(ctx)
	if err != nil {
		return res, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer c.Release()

	t, err := c.BeginTx(ctx, txOptions)
	if err != nil {
		return res, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer t.Rollback(ctx)

	res, err = fn(&Context{
		UserRepository: repositories.NewUserRepository(t.Conn()),
	})
	if err != nil {
		return res, fmt.Errorf("failed to execute transaction: %w", err)
	}

	err = t.Commit(ctx)
	if err != nil {
		return res, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return res, nil
}
