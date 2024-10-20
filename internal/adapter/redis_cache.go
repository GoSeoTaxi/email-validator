package adapter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/GoSeoTaxi/email-validator/internal/domain"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) domain.Cache {
	return &RedisCache{client: client}
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}
