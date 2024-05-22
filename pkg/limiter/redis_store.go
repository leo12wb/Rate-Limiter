package limiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

var _ CounterStore = &RedisStore{}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client}
}

func (r *RedisStore) Increment(key string) (int64, error) {
	ctx := context.Background()
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if err := r.client.ExpireNX(ctx, key, time.Second).Err(); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RedisStore) Set(key string, value int64, timeout time.Duration) error {
	ctx := context.Background()
	if err := r.client.Set(ctx, key, value, timeout).Err(); err != nil {
		return err
	}
	return nil
}
