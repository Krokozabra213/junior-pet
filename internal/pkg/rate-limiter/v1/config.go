package ratelimiterv1

import (
	"sync"
	"time"
)

// Limit описывает ограничение: количество запросов за интервал времени.
type Limit struct {
	count  int
	window time.Duration
}

func NewLimit(count int, window time.Duration) Limit {
	return Limit{
		count:  count,
		window: window,
	}
}

// Config хранит настройки ограничений для методов.
type Config struct {
	defaultLimit Limit

	// Лимиты для конкретных gRPC методов
	// Ключ: полное имя метода, например "/sso.AuthService/Login"
	methodLimits map[string]Limit
	mu           *sync.RWMutex
}

// NewConfig создаёт конфигурацию с указанным лимитом по умолчанию.
func NewConfig(count int, window time.Duration) *Config {

	if count <= 0 || window <= 0 {
		panic("invalid count or window: must be positive")
	}

	return &Config{
		defaultLimit: NewLimit(count, window),
		methodLimits: make(map[string]Limit),
		mu:           &sync.RWMutex{},
	}
}

// SetMethod устанавливает лимит для конкретного метода.
// Если метод уже был настроен, лимит перезаписывается.
func (c *Config) SetMethod(method string, count int, window time.Duration) {
	if count <= 0 || window <= 0 {
		panic("invalid limit: count and window must be positive")
	}

	c.mu.Lock()
	c.methodLimits[method] = NewLimit(count, window)
	c.mu.Unlock()
}

func (c *Config) getMethodLimitOrDefault(method string) Limit {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if ml, ok := c.methodLimits[method]; ok {
		return ml
	}
	return c.defaultLimit
}
