//go:build integration

package tests

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	repository "github.com/Krokozabra213/schools_backend/services/sso/repository/postgres"
	"github.com/Krokozabra213/schools_backend/sql/goose/sso/migrations"
)

var (
	testPool *pgxpool.Pool
	testRepo *repository.PostgresRepository
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Запускаем PostgreSQL
	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		fmt.Printf("failed to start container: %v\n", err)
		os.Exit(1)
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Printf("failed to get connection string: %v\n", err)
		os.Exit(1)
	}

	testPool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		fmt.Printf("failed to connect: %v\n", err)
		os.Exit(1)
	}

	// Запускаем goose миграции
	if err := runMigrations(ctx, connStr); err != nil {
		fmt.Printf("failed to migrate: %v\n", err)
		os.Exit(1)
	}

	testRepo = repository.NewRepository(testPool)

	code := m.Run()

	testPool.Close()
	container.Terminate(ctx)
	os.Exit(code)
}

func runMigrations(ctx context.Context, connStr string) error {
	// Goose работает с database/sql
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	// Используем embed.FS с миграциями
	goose.SetBaseFS(migrations.Files)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}

func cleanup(t *testing.T) {
	t.Helper()
	_, err := testPool.Exec(context.Background(), "TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}
}
