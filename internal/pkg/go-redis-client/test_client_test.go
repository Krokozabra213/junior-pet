package gorediscli

import (
	"context"
	"testing"
	"time"
)

func TestNewFailsWithoutPassword(t *testing.T) {
	ctx := context.Background()

	_, err := New(ctx,
		WithAddr("localhost:6379"),
		// password не передан
	)

	if err == nil {
		t.Error("expected error when password is not set")
	}
}

func TestNewFailsWithInvalidConfig(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		opts []Option
	}{
		{
			name: "no password",
			opts: []Option{WithAddr("localhost:6379")},
		},
		{
			name: "empty addr",
			opts: []Option{WithAddr(""), WithPassword("test")},
		},
		{
			name: "invalid db",
			opts: []Option{WithPassword("test"), WithDB(16)},
		},
		{
			name: "zero pool",
			opts: []Option{WithPassword("test"), WithPoolSize(0)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(ctx, tt.opts...)
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

func TestNewFailsOnUnreachableRedis(t *testing.T) {
	ctx := context.Background()

	_, err := New(ctx,
		WithAddr("localhost:59999"), // порт где ничего нет
		WithPassword("test"),
		WithPingTimeout(500*time.Millisecond),
		WithDialTimeout(500*time.Millisecond),
	)

	if err == nil {
		t.Error("expected error when redis is unreachable")
	}
}
