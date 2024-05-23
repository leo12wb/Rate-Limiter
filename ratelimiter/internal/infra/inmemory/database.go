package inmemory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/leo12wb/Rate-Limiter/ratelimiter/internal/entity"
)

type repository struct {
	limiters map[string]*entity.RateLimiterInfo
	mu       sync.Mutex
}

func NewDatabaseRepository() entity.DatabaseRepository {
	return &repository{
		limiters: make(map[string]*entity.RateLimiterInfo),
	}
}

func (r *repository) Create(ctx context.Context, config entity.RateLimiterInfo, window time.Duration) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := fmt.Sprintf("ratelimitconfig:%s", config.Key)
	r.limiters[key] = &config
	return nil
}

func (r *repository) Read(ctx context.Context, key string) (*entity.RateLimiterInfo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key = fmt.Sprintf("ratelimitconfig:%s", key)
	if val, exists := r.limiters[key]; exists {
		return val, nil
	}
	return nil, fmt.Errorf("rate limiter not found")
}

func (r *repository) CheckLimit(ctx context.Context, key string, requests int, window time.Duration) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key = fmt.Sprintf("ratelimit:%s", key)
	limiter, exists := r.limiters[key]
	if !exists {
		limiter = &entity.RateLimiterInfo{
			Key:       key,
			Requests:  requests,
			Remaining: requests,
			Reset:     time.Now().Add(window).Unix(),
		}
		r.limiters[key] = limiter
	}

	if time.Now().Unix() > limiter.Reset {
		limiter.Remaining = requests
		limiter.Reset = time.Now().Add(window).Unix()
	}

	if limiter.Remaining > 0 {
		limiter.Remaining--
		return true, nil
	}

	return false, nil
}
