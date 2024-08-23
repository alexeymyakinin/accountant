package repositories

import (
	"accountant/internal/domain"
	"accountant/internal/infra/db"
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type AccountRepository struct {
	conn *pgx.Conn
}

func (r *AccountRepository) ListByUserId(ctx context.Context, userId int) ([]domain.Account, error) {
	query, args := squirrel.Select(
		"id",
		"name",
		"type",
	).From(
		"accounts",
	).Where(
		squirrel.Eq{"user_id": userId},
	).PlaceholderFormat(
		squirrel.Dollar,
	).MustSql()

	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[db.Account])
	if err != nil {
		return nil, err
	}

	res := make([]domain.Account, len(accounts))
	for i, row := range accounts {
		res[i] = domain.Account{
			Id:   row.Id,
			Name: row.Name,
			Type: row.Type,
		}
	}

	return res, nil
}
