//go:build integration

package gorediscli_test

import (
	"context"
	"testing"
	"time"

	gorediscli "github.com/Krokozabra213/schools_backend/internal/pkg/go-redis-client"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupRedis(t *testing.T, ctx context.Context) (string, func()) {
	t.Helper()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:7-alpine",
			ExposedPorts: []string{"6379/tcp"},
			Cmd:          []string{"redis-server", "--requirepass", "testpass"},
			WaitingFor:   wait.ForLog("Ready to accept connections").WithStartupTimeout(30 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "6379")
	addr := host + ":" + port.Port()

	cleanup := func() {
		container.Terminate(ctx)
	}

	return addr, cleanup
}

func TestNewIntegration(t *testing.T) {
	ctx := context.Background()
	addr, cleanup := setupRedis(t, ctx)
	defer cleanup()

	client, err := gorediscli.New(ctx,
		gorediscli.WithAddr(addr),
		gorediscli.WithPassword("testpass"),
		gorediscli.WithDB(0),
	)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer client.Close()

	// Ping
	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatalf("Ping error: %v", err)
	}
}

func TestSetGetIntegration(t *testing.T) {
	ctx := context.Background()
	addr, cleanup := setupRedis(t, ctx)
	defer cleanup()

	client, err := gorediscli.New(ctx,
		gorediscli.WithAddr(addr),
		gorediscli.WithPassword("testpass"),
	)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer client.Close()

	// Set
	err = client.Set(ctx, "testkey", "testvalue", 10*time.Second).Err()
	if err != nil {
		t.Fatalf("Set error: %v", err)
	}

	// Get
	val, err := client.Get(ctx, "testkey").Result()
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if val != "testvalue" {
		t.Errorf("got %q, want %q", val, "testvalue")
	}

	// Del
	deleted, err := client.Del(ctx, "testkey").Result()
	if err != nil {
		t.Fatalf("Del error: %v", err)
	}
	if deleted != 1 {
		t.Errorf("deleted %d keys, want 1", deleted)
	}
}

func TestWrongPasswordIntegration(t *testing.T) {
	ctx := context.Background()
	addr, cleanup := setupRedis(t, ctx)
	defer cleanup()

	_, err := gorediscli.New(ctx,
		gorediscli.WithAddr(addr),
		gorediscli.WithPassword("wrongpassword"),
		gorediscli.WithPingTimeout(2*time.Second),
	)

	if err == nil {
		t.Error("expected error with wrong password")
	}
}
