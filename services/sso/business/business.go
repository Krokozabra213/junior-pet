package business

import (
	"context"
	"log/slog"

	ssoconfig "github.com/Krokozabra213/schools_backend/services/sso/config"
	"github.com/Krokozabra213/schools_backend/services/sso/domain"
)

//go:generate mockgen  -source=constructor.go -destination=mocks/mocks.go

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

type Business struct {
	cfg  *ssoconfig.Config
    log *slog.Logger
	user UserProvider
}

func New(cfg *ssoconfig.Config, user UserProvider) *Business {
	return &Business{
		cfg:  cfg,
		user: user,
	}
}

func (b *Business) checkUpdatePermission(ctx context.Context, actorID, targetID int64) error {
	// Пользователь может редактировать только себя
	if actorID != targetID {
		return ErrPermissionDenied
	}
	return nil
}
