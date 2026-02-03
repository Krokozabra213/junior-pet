//go:build integration

package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Krokozabra213/schools_backend/services/sso/domain"
	repository "github.com/Krokozabra213/schools_backend/services/sso/repository/postgres"
)

func TestCreateUser(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		cleanup(t)

		result, err := testRepo.CreateUser(ctx, &domain.CreateUser{
			Username: "john",
			Email:    "john@test.com",
			Password: "password123",
			Name:     "John",
			Surname:  "Doe",
			IsMale:   true,
		})

		require.NoError(t, err)
		assert.NotZero(t, result.ID)
		assert.NotZero(t, result.CreatedAt)
	})

	t.Run("duplicate username", func(t *testing.T) {
		cleanup(t)

		user := &domain.CreateUser{
			Username: "duplicate",
			Email:    "first@test.com",
			Password: "password",
			Name:     "John",
			Surname:  "Doe",
			IsMale:   true,
		}

		_, err := testRepo.CreateUser(ctx, user)
		require.NoError(t, err)

		user.Email = "second@test.com"
		_, err = testRepo.CreateUser(ctx, user)

		assert.ErrorIs(t, err, repository.ErrAlreadyExists)
	})

	t.Run("duplicate email", func(t *testing.T) {
		cleanup(t)

		user := &domain.CreateUser{
			Username: "user1",
			Email:    "same@test.com",
			Password: "password",
			Name:     "John",
			Surname:  "Doe",
			IsMale:   true,
		}

		_, err := testRepo.CreateUser(ctx, user)
		require.NoError(t, err)

		user.Username = "user2"
		_, err = testRepo.CreateUser(ctx, user)

		assert.ErrorIs(t, err, repository.ErrAlreadyExists)
	})
}

func TestGetUserByID(t *testing.T) {
	ctx := context.Background()

	t.Run("found", func(t *testing.T) {
		cleanup(t)

		created, _ := testRepo.CreateUser(ctx, &domain.CreateUser{
			Username: "findme",
			Email:    "find@test.com",
			Password: "password",
			Name:     "John",
			Surname:  "Doe",
			IsMale:   true,
		})

		user, err := testRepo.GetUserByID(ctx, created.ID)

		require.NoError(t, err)
		assert.Equal(t, "findme", user.Username)
		assert.Equal(t, "find@test.com", user.Email)
		assert.Equal(t, "John", user.Name)
	})

	t.Run("not found", func(t *testing.T) {
		cleanup(t)

		_, err := testRepo.GetUserByID(ctx, 99999)

		assert.ErrorIs(t, err, repository.ErrNotFound)
	})
}

func TestGetUserByUsername(t *testing.T) {
	ctx := context.Background()

	t.Run("found", func(t *testing.T) {
		cleanup(t)

		testRepo.CreateUser(ctx, &domain.CreateUser{
			Username: "searchme",
			Email:    "search@test.com",
			Password: "password",
			Name:     "Jane",
			Surname:  "Doe",
			IsMale:   false,
		})

		user, err := testRepo.GetUserByUsername(ctx, "searchme")

		require.NoError(t, err)
		assert.Equal(t, "search@test.com", user.Email)
	})

	t.Run("not found", func(t *testing.T) {
		cleanup(t)

		_, err := testRepo.GetUserByUsername(ctx, "nonexistent")

		assert.ErrorIs(t, err, repository.ErrNotFound)
	})
}

func TestUpdatePassword(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		cleanup(t)

		created, _ := testRepo.CreateUser(ctx, &domain.CreateUser{
			Username: "user",
			Email:    "user@test.com",
			Password: "old_password",
			Name:     "John",
			Surname:  "Doe",
			IsMale:   true,
		})

		err := testRepo.UpdatePassword(ctx, created.ID, "new_password")
		require.NoError(t, err)

		user, _ := testRepo.GetUserByID(ctx, created.ID)
		assert.Equal(t, "new_password", user.Password)
	})
}

func TestSoftDeleteUser(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		cleanup(t)

		created, _ := testRepo.CreateUser(ctx, &domain.CreateUser{
			Username: "todelete",
			Email:    "delete@test.com",
			Password: "password",
			Name:     "John",
			Surname:  "Doe",
			IsMale:   true,
		})

		err := testRepo.SoftDeleteUser(ctx, created.ID)
		require.NoError(t, err)

		_, err = testRepo.GetUserByID(ctx, created.ID)
		assert.ErrorIs(t, err, repository.ErrNotFound)
	})
}

func TestExistsUserByUsername(t *testing.T) {
	ctx := context.Background()

	t.Run("exists", func(t *testing.T) {
		cleanup(t)

		testRepo.CreateUser(ctx, &domain.CreateUser{
			Username: "exists",
			Email:    "exists@test.com",
			Password: "password",
			Name:     "John",
			Surname:  "Doe",
			IsMale:   true,
		})

		exists, err := testRepo.ExistsUserByUsername(ctx, "exists")

		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("not exists", func(t *testing.T) {
		cleanup(t)

		exists, err := testRepo.ExistsUserByUsername(ctx, "nonexistent")

		require.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestExistsUserByEmail(t *testing.T) {
	ctx := context.Background()

	t.Run("exists", func(t *testing.T) {
		cleanup(t)

		testRepo.CreateUser(ctx, &domain.CreateUser{
			Username: "user",
			Email:    "check@test.com",
			Password: "password",
			Name:     "John",
			Surname:  "Doe",
			IsMale:   true,
		})

		exists, err := testRepo.ExistsUserByEmail(ctx, "check@test.com")

		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("not exists", func(t *testing.T) {
		cleanup(t)

		exists, err := testRepo.ExistsUserByEmail(ctx, "nonexistent@test.com")

		require.NoError(t, err)
		assert.False(t, exists)
	})
}
