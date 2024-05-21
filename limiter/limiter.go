package limiter

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RateLimiter struct {
	client *redis.Client
	limiterConfig map[string]int
	ctx context.Context
}

func NewRateLimiter() (*RateLimiter, error) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	limiterConfig := make(map[string]int)

	ipLimit, err := strconv.Atoi(os.Getenv("IP_LIMIT"))
	if err != nil {
		return nil, err
	}
	tokenLimit, err := strconv.Atoi(os.Getenv("TOKEN_LIMIT"))
	if err != nil {
		return nil, err
	}

	limiterConfig["ip"] = ipLimit
	limiterConfig["token"] = tokenLimit

	return &RateLimiter{
		client: client,
		limiterConfig: limiterConfig,
		ctx: ctx,
	}, nil
}

func (r *RateLimiter) Allow(identifier string, maxRequests int, window time.Duration) (bool, error) {
	allowed, err := r.client.SetNX(r.ctx, identifier, 0, window).Result()
	if err != nil {
		return false, err
	}
	if allowed {
		r.client.Expire(r.ctx, identifier, window)
	}
	reqCount, err := r.client.Incr(r.ctx, identifier).Result()
	if err != nil {
		return false, err
	}
	if reqCount > int64(maxRequests) {
		return false, nil
	}
	return true, nil
}

func (r *RateLimiter) Get(identifier string) (int, error) {
	return 0, nil
}

func (r *RateLimiter) Incr(identifier string) error {
	return nil
}

func (r *RateLimiter) Decr(identifier string) error {
	return nil
}
