package business

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Krokozabra213/schools_backend/services/sso/domain"
	"github.com/Krokozabra213/schools_backend/services/sso/repository/postgres"
)

func (b *Business) UpdateUser(ctx context.Context, actorID int64, params domain.UpdateUser) error {
	const op = "business.UpdateUser"

	log := b.log.With(
		slog.String("op", op),
		slog.Int64("target_user_id", params.ID),
		slog.Int64("actor_id", actorID),
	)
	log.Info("starting update user process...")

	// 1. Проверка прав доступа
	if err := b.checkUpdatePermission(ctx, actorID, params.ID); err != nil {
		log.Warn("permission denied", slog.String("error", err.Error()))
		return ErrPermissionDenied
	}

	err := b.user.UpdateUser(ctx, params)
	if err != nil {
		log.Error("failed update user", slog.String("error", err.Error()))
		if errors.Is(err, postgres.ErrNotFound) {
			return ErrUserNotFound
		}
		return ErrInternal
	}

	log.Info("user successfully updated")
	return nil
}
