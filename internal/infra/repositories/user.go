package repositories

import (
	"accountant/internal/domain"
	"accountant/internal/infra/db"
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	conn *pgx.Conn
}

func NewUserRepository(conn *pgx.Conn) *UserRepository {
	return &UserRepository{conn: conn}
}

func (repo *UserRepository) Insert(ctx context.Context, email string, hashedPassword string) (int, error) {
	query, queryArgs := squirrel.Insert(
		"users",
	).Columns(
		"email",
		"password",
	).Values(
		email,
		hashedPassword,
	).Suffix(
		"RETURNING \"id\"",
	).PlaceholderFormat(
		squirrel.Dollar,
	).MustSql()

	rows, err := repo.conn.Query(ctx, query, queryArgs...)
	if err != nil {
		return 0, err
	}

	id, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *UserRepository) GetById(ctx context.Context, id int) (domain.User, error) {
	query, queryArgs := squirrel.Select(
		"*",
	).From(
		"users",
	).Where(
		squirrel.Eq{"id": id},
	).PlaceholderFormat(
		squirrel.Dollar,
	).MustSql()

	rows, err := repo.conn.Query(ctx, query, queryArgs...)
	if err != nil {
		return domain.User{}, err
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[db.User])
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return domain.User{}, NotFoundError
		default:
			return domain.User{}, err
		}
	}

	return domain.User{
		Id:             user.Id,
		Email:          user.Email,
		HashedPassword: user.Password,
	}, nil
}

func (repo *UserRepository) IsEmailExists(ctx context.Context, email string) (bool, error) {
	sql, args := squirrel.Select(
		"email",
	).From(
		"users",
	).Where(
		squirrel.Eq{"email": email},
	).PlaceholderFormat(
		squirrel.Dollar,
	).MustSql()

	r, err := repo.conn.Query(ctx, sql, args...)
	if err != nil {
		return false, nil
	}
	defer r.Close()

	count := 0
	for r.Next() {
		count++
	}

	return count > 0, nil
}

func (repo *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	query, queryArgs := squirrel.Select(
		"*",
	).From(
		"users",
	).Where(
		squirrel.Eq{"email": email},
	).PlaceholderFormat(
		squirrel.Dollar,
	).MustSql()

	rows, err := repo.conn.Query(ctx, query, queryArgs...)
	if err != nil {
		return domain.User{}, err
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[db.User])
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return domain.User{}, NotFoundError
		default:
			return domain.User{}, err
		}
	}

	return domain.User{
		Id:             user.Id,
		Email:          user.Email,
		HashedPassword: user.Password,
	}, nil
}
