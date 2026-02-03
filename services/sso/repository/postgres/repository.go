package postgres

import (
	"context"
	"errors"

	"github.com/Krokozabra213/schools_backend/services/sso/domain"
	"github.com/Krokozabra213/schools_backend/services/sso/repository/postgres/sqlc"
	"github.com/jackc/pgx/v5"
)

type UserProvider interface {
	UpdateUser(ctx context.Context, params domain.UpdateUser) error
	CreateUser(ctx context.Context, user *domain.CreateUser) (*domain.CreateUserRow, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)
	UpdatePassword(ctx context.Context, id int64, password string) error
	SoftDeleteUser(ctx context.Context, id int64) error
	HardDeleteUser(ctx context.Context, id int64) error
	CountUsers(ctx context.Context) (int64, error)
	ExistsUserByUsername(ctx context.Context, username string) (bool, error)
	ExistsUserByEmail(ctx context.Context, email string) (bool, error)
}

var _ UserProvider = (*PostgresRepository)(nil)

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
