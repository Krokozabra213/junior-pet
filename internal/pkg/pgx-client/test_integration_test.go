//go:build integration

package pgxclient_test

import (
	"context"
	"testing"
	"time"

	pgxclient "github.com/Krokozabra213/schools_backend/internal/pkg/pgx-client"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestNewIntegration(t *testing.T) {
	ctx := context.Background()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:16-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     "test",
				"POSTGRES_PASSWORD": "test",
				"POSTGRES_DB":       "testdb",
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(ctx)

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	client, err := pgxclient.New(ctx,
		pgxclient.WithHost(host),
		pgxclient.WithPort(port.Int()),
		pgxclient.WithUser("test"),
		pgxclient.WithPassword("test"),
		pgxclient.WithDatabase("testdb"),
	)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer client.Close()

	// Проверяем что реально работает
	var result int
	err = client.QueryRow(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		t.Fatalf("QueryRow error: %v", err)
	}
	if result != 1 {
		t.Errorf("got %d, want 1", result)
	}
}
