package rdb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/leo12wb/Rate-Limiter/ratelimiter/internal/entity"

	"github.com/go-redis/redis/v8"
)

type repository struct {
	database *redis.Client
}

func NewDatabaseRepository(database *redis.Client) entity.DatabaseRepository {
	return &repository{database: database}
}

func (repository repository) Create(ctx context.Context, config entity.RateLimiterInfo, every time.Duration) error {
	key := fmt.Sprintf("ratelimitconfig:%s", config.Key)
	err := repository.database.Set(ctx, key, config, every).Err()
	if err != nil {
		return err
	}

	return nil
}

func (repository repository) Read(ctx context.Context, key string) (*entity.RateLimiterInfo, error) {
	key = fmt.Sprintf("ratelimitconfig:%s", key)
	val, err := repository.database.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var config entity.RateLimiterInfo
	err = json.Unmarshal([]byte(val), &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (repository repository) CheckLimit(ctx context.Context, key string, requests int, every time.Duration) (bool, error) {
	key = fmt.Sprintf("ratelimit:%s", key)
	current, err := repository.database.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return true, err
	}

	if current >= requests {
		return false, nil
	}

	_, err = repository.database.Incr(ctx, key).Result()
	if err != nil {
		return true, err
	}

	if current == 0 {
		_, err = repository.database.Expire(ctx, key, every).Result()
		if err != nil {
			return true, err
		}
	}

	return true, nil
}
