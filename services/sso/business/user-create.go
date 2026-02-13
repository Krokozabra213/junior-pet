package business

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Krokozabra213/schools_backend/services/sso/domain"
	"github.com/Krokozabra213/schools_backend/services/sso/repository/postgres"
	"golang.org/x/crypto/bcrypt"
)

func (b *Business) CreateUser(ctx context.Context, user *domain.CreateUser) (*domain.CreateUserRow, error) {
	const op = "business.CreateUser"

	log := b.log.With(
		slog.String("op", op),
		slog.String("username", user.Username),
	)
	log.Info("starting user registration process...")

	exists, err := b.user.ExistsUserByEmail(ctx, user.Email)
	if err != nil {
		log.Error("failed to check email existence", slog.String("error", err.Error()))
		return nil, ErrInternal
	}
	if exists {
		log.Warn("email already exists")
		return nil, ErrEmailExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", slog.String("error", err.Error()))
		return nil, err
	}

	// если пароль потребуется дальше
	_ = user.Password
	user.Password = string(hash)

	result, err := b.user.CreateUser(ctx, user)
	if err != nil {
		log.Error("failed create new user", slog.String("error", err.Error()))
		if errors.Is(err, postgres.ErrAlreadyExists) {
			return nil, ErrUserExists
		}
		return nil, ErrInternal
	}

	log.Info("user successfully registered",
		slog.Int64("user_id", result.ID))

	return result, nil
}
