//go:build integration

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Krokozabra213/schools_backend/services/sso/domain"
	repository "github.com/Krokozabra213/schools_backend/services/sso/repository/redis"
)

func TestCacheUserProfile(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		cleanup(t)

		profile := &domain.UserCacheProfile{
			ID:       123,
			Username: "testuser",
			Email:    "test@example.com",
			Name:     "John",
			Surname:  "Doe",
			IsMale:   true,
		}

		err := testRepo.CacheUserProfile(ctx, profile, 10*time.Minute)

		require.NoError(t, err)

		// Проверяем что сохранилось
		cached, err := testRepo.GetUserProfile(ctx, 123)
		require.NoError(t, err)
		assert.Equal(t, profile.Username, cached.Username)
		assert.Equal(t, profile.Email, cached.Email)
		assert.Equal(t, profile.Name, cached.Name)
		assert.Equal(t, profile.Surname, cached.Surname)
		assert.Equal(t, profile.IsMale, cached.IsMale)
	})

	t.Run("overwrites existing", func(t *testing.T) {
		cleanup(t)

		// Сохраняем первый профиль
		profile1 := &domain.UserCacheProfile{
			ID:       123,
			Username: "olduser",
			Email:    "old@example.com",
		}
		testRepo.CacheUserProfile(ctx, profile1, 10*time.Minute)

		// Перезаписываем
		profile2 := &domain.UserCacheProfile{
			ID:       123,
			Username: "newuser",
			Email:    "new@example.com",
		}
		err := testRepo.CacheUserProfile(ctx, profile2, 10*time.Minute)
		require.NoError(t, err)

		// Проверяем что обновилось
		cached, _ := testRepo.GetUserProfile(ctx, 123)
		assert.Equal(t, "newuser", cached.Username)
		assert.Equal(t, "new@example.com", cached.Email)
	})
}

func TestGetUserProfile(t *testing.T) {
	ctx := context.Background()

	t.Run("found", func(t *testing.T) {
		cleanup(t)

		profile := &domain.UserCacheProfile{
			ID:       456,
			Username: "findme",
			Email:    "find@example.com",
			Name:     "Jane",
			Surname:  "Smith",
			IsMale:   false,
		}
		testRepo.CacheUserProfile(ctx, profile, 10*time.Minute)

		cached, err := testRepo.GetUserProfile(ctx, 456)

		require.NoError(t, err)
		assert.Equal(t, profile.ID, cached.ID)
		assert.Equal(t, profile.Username, cached.Username)
		assert.Equal(t, profile.Email, cached.Email)
	})

	t.Run("not found", func(t *testing.T) {
		cleanup(t)

		cached, err := testRepo.GetUserProfile(ctx, 99999)

		assert.Nil(t, cached)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})

	t.Run("expired", func(t *testing.T) {
		cleanup(t)

		profile := &domain.UserCacheProfile{
			ID:       789,
			Username: "expiring",
		}
		testRepo.CacheUserProfile(ctx, profile, 2*time.Second)

		// Сначала есть
		cached, err := testRepo.GetUserProfile(ctx, 789)
		require.NoError(t, err)
		assert.NotNil(t, cached)

		// Ждём истечения
		time.Sleep(3 * time.Second)

		// После TTL - нет
		cached, err = testRepo.GetUserProfile(ctx, 789)
		assert.Nil(t, cached)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})
}

func TestDeleteUserProfile(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		cleanup(t)

		profile := &domain.UserCacheProfile{
			ID:       123,
			Username: "todelete",
		}
		testRepo.CacheUserProfile(ctx, profile, 10*time.Minute)

		// Удаляем
		err := testRepo.DeleteUserProfile(ctx, 123)
		require.NoError(t, err)

		// Проверяем что удалено
		cached, err := testRepo.GetUserProfile(ctx, 123)
		assert.Nil(t, cached)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})

	t.Run("not exists - no error", func(t *testing.T) {
		cleanup(t)

		// Удаляем несуществующий - не должно быть ошибки
		err := testRepo.DeleteUserProfile(ctx, 99999)

		assert.NoError(t, err)
	})
}

func TestUpdateUserProfileTTL(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		cleanup(t)

		profile := &domain.UserCacheProfile{
			ID:       123,
			Username: "ttltest",
		}
		testRepo.CacheUserProfile(ctx, profile, 2*time.Second)

		// Продляем TTL
		err := testRepo.UpdateUserProfileTTL(ctx, 123, 10*time.Second)
		require.NoError(t, err)

		// Ждём старый TTL
		time.Sleep(3 * time.Second)

		// Всё ещё существует (потому что продлили)
		cached, err := testRepo.GetUserProfile(ctx, 123)
		require.NoError(t, err)
		assert.NotNil(t, cached)
	})

	t.Run("key not exists", func(t *testing.T) {
		cleanup(t)

		// Пытаемся продлить несуществующий ключ
		err := testRepo.UpdateUserProfileTTL(ctx, 99999, 10*time.Second)
		t.Log(err)

		// Зависит от реализации - может быть ошибка или нет
		// Redis EXPIRE возвращает 0 если ключа нет
		assert.NoError(t, err) // или assert.Error если ты возвращаешь ошибку
	})
}
