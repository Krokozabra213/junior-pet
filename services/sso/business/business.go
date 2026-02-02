package business

import (
	"context"

	"github.com/Krokozabra213/schools_backend/services/sso/domain"
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

type Business struct {
	cfg          *ssonewconfig.Config
	tokenRepo    ITokenProvider
	userProvider IUserProvider
	appProvider  IAppProvider
	jwtManager   IJWTManager
	hasher       IHasher
	publicKeyPEM string
}
