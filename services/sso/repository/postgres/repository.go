package postgres

import (
	"context"
	"errors"

	"github.com/Krokozabra213/schools_backend/services/sso/business"
	"github.com/Krokozabra213/schools_backend/services/sso/repository/postgres/sqlc"
	"github.com/jackc/pgx/v5"
)

var _ business.UserProvider = (*PostgresRepository)(nil)

type PostgresRepository struct {
	DB      sqlc.DBTX
	Queries sqlc.Querier
}

func NewRepository(db sqlc.DBTX) *PostgresRepository {
	return &PostgresRepository{
		DB:      db,
		Queries: sqlc.New(db),
	}
}

func (r *PostgresRepository) WithTx(tx pgx.Tx) *PostgresRepository {
	return &PostgresRepository{
		DB:      tx,
		Queries: sqlc.New(tx),
	}
}

func (r *PostgresRepository) handleError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return err
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return ErrInternal
}
