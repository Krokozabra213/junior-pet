package postgres

import (
	"context"
	"errors"

	"github.com/Krokozabra213/schools_backend/services/sso/domain"
	"github.com/Krokozabra213/schools_backend/services/sso/repository/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *PostgresRepository) CreateUser(ctx context.Context, user *domain.CreateUser) (*domain.CreateUserRow, error) {
	result, err := r.Queries.CreateUser(ctx, sqlc.CreateUserParams{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Name:     user.Name,
		Surname:  user.Surname,
		IsMale:   user.IsMale,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrUserAlreadyExists
		}
		return nil, r.handleError(err)
	}

	return &domain.CreateUserRow{
		ID:        result.ID,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

func (r *PostgresRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	result, err := r.Queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, r.handleError(err)
	}
	user := domain.NewUser(result.ID, result.Username, result.Email, result.Password, result.Name,
		result.Surname, result.IsMale, result.CreatedAt, result.UpdatedAt)

	return &user, nil
}

func (r *PostgresRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	result, err := r.Queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, r.handleError(err)
	}
	user := domain.NewUser(result.ID, result.Username, result.Email, result.Password, result.Name,
		result.Surname, result.IsMale, result.CreatedAt, result.UpdatedAt)

	return &user, nil
}

func (r *PostgresRepository) UpdatePassword(ctx context.Context, id int64, password string) error {
	err := r.Queries.UpdatePassword(ctx, sqlc.UpdatePasswordParams{
		ID:       id,
		Password: password,
	})
	if err != nil {
		return r.handleError(err)
	}
	return nil
}

func (r *PostgresRepository) SoftDeleteUser(ctx context.Context, id int64) error {
	err := r.Queries.SoftDeleteUser(ctx, id)
	if err != nil {
		return r.handleError(err)
	}
	return nil
}

func (r *PostgresRepository) HardDeleteUser(ctx context.Context, id int64) error {
	err := r.Queries.HardDeleteUser(ctx, id)
	if err != nil {
		return r.handleError(err)
	}
	return nil
}

func (r *PostgresRepository) CountUsers(ctx context.Context) (int64, error) {
	count, err := r.Queries.CountUsers(ctx)
	if err != nil {
		return 0, r.handleError(err)
	}
	return count, nil
}

func (r *PostgresRepository) ExistsUserByUsername(ctx context.Context, username string) (bool, error) {
	exists, err := r.Queries.ExistsUserByUsername(ctx, username)
	if err != nil {
		return false, r.handleError(err)
	}

	return exists, nil
}

func (r *PostgresRepository) ExistsUserByEmail(ctx context.Context, email string) (bool, error) {
	exists, err := r.Queries.ExistsUserByEmail(ctx, email)
	if err != nil {
		return false, r.handleError(err)
	}

	return exists, nil
}
