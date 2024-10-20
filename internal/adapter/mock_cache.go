package adapter

import (
	"context"
	"sync"
	"time"

	"github.com/GoSeoTaxi/email-validator/internal/domain"
)

type MockCache struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewMockCache() domain.Cache {
	return &MockCache{
		data: make(map[string]string),
	}
}

func (m *MockCache) Get(ctx context.Context, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, exists := m.data[key]
	if !exists {
		return "", domain.ErrCacheMiss
	}
	return val, nil
}

func (m *MockCache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
	return nil
}
