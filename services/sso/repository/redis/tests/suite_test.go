//go:build integration

package tests

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"

	redisclient "github.com/Krokozabra213/schools_backend/pkg/redis-client"
	repository "github.com/Krokozabra213/schools_backend/services/sso/repository/redis"
)

var (
	testClient *redisclient.RedisClient
	testRepo   *repository.RedisRepository
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Запускаем Redis
	container, err := redis.Run(ctx,
		"redis:7-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		fmt.Printf("failed to start container: %v\n", err)
		os.Exit(1)
	}

	// Получаем connection string
	connStr, err := container.ConnectionString(ctx)
	if err != nil {
		fmt.Printf("failed to get connection string: %v\n", err)
		os.Exit(1)
	}

	addr := strings.TrimPrefix(connStr, "redis://")

	// Создаём клиент
	cfg := redisclient.NewConfig(
		addr,           // addr (формат host:port)
		"",             // password
		0,              // db
		5*time.Second,  // dial timeout
		3*time.Second,  // read timeout
		3*time.Second,  // write timeout
		1*time.Hour,    // max conn lifetime
		30*time.Minute, // max conn idle time
		10,             // pool size
		2,              // min idle conns
	)

	testClient, err = redisclient.New(ctx, cfg)
	if err != nil {
		fmt.Printf("failed to connect: %v\n", err)
		os.Exit(1)
	}

	testRepo = repository.NewRepository(testClient)

	code := m.Run()

	testClient.CloseConn()
	container.Terminate(ctx)
	os.Exit(code)
}

// cleanup очищает все ключи между тестами
func cleanup(t *testing.T) {
	t.Helper()
	ctx := context.Background()

	if err := testClient.FlushDB(ctx).Err(); err != nil {
		t.Fatalf("failed to cleanup: %v", err)
	}
}
