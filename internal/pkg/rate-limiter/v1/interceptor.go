package ratelimiterv1

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const limiterPrefix = "ratelimit"

// UnaryInterceptor returns a gRPC unary interceptor for rate limiting.
func UnaryInterceptor(limiter Limiter, cfg Config, log *slog.Logger) grpc.UnaryServerInterceptor {

	if log == nil {
		log = slog.Default()
	}

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		clientIP, err := extractClientIP(ctx)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed get ip address")
		}
		method := info.FullMethod

		ml := cfg.getMethodLimitOrDefault(method)

		// Ключ: ip + method
		key := rateLimitKey(clientIP, method)

		allowed, err := limiter.Allow(ctx, key, ml.count, ml.window)
		if err != nil {
			// При ошибке Redis — пропускаем (fail open)
			// В проде можно сделать fail closed
			log.Error("rate limiter error",
				slog.String("error", err.Error()),
				slog.String("ip", clientIP),
				slog.String("method", method),
			)
			return handler(ctx, req)
		}

		if !allowed {
			log.Warn("rate limit exceeded",
				slog.String("ip", clientIP),
				slog.String("method", method),
				slog.Int("limit", ml.count),
				slog.Duration("window", ml.window),
			)

			return nil, status.Error(codes.ResourceExhausted,
				"rate limit exceeded, try again later",
			)
		}

		return handler(ctx, req)
	}
}

// extractClientIP gets client IP from gRPC context.
// Checks X-Forwarded-For first (behind load balancer), then peer address.
func extractClientIP(ctx context.Context) (string, error) {
	// 1. Proxy/LB header
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if xff := md.Get("x-forwarded-for"); len(xff) > 0 {
			ips := strings.SplitN(xff[0], ",", 2)
			return strings.TrimSpace(ips[0]), nil
		}
		if xri := md.Get("x-real-ip"); len(xri) > 0 {
			return xri[0], nil
		}
	}

	// 2. Direct peer address
	if p, ok := peer.FromContext(ctx); ok {
		host, _, err := net.SplitHostPort(p.Addr.String())
		if err == nil {
			return host, nil
		}
		return p.Addr.String(), nil // fallback
	}

	return "unknown", errors.New("bad ip")
}

func rateLimitKey(ip, method string) string {
	// Заменяем двоеточия в IPv6 на подчёркивания, чтобы не ломать разделитель
	safeIP := strings.ReplaceAll(ip, ":", "_")
	return limiterPrefix + safeIP + "#" + method
}
