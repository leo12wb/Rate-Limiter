package rate_limit

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/leo12wb/Rate-Limiter/internal/entity/web_session"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RateLimitRepositoryRedis struct {
	Client *redis.Client // Unique to redis, can't use sql.DB abstraction here.
}

func NewRateLimitRepositoryRedis(client *redis.Client) *RateLimitRepository {
	var rateLimiter RateLimitRepository
	rateLimiter = &RateLimitRepositoryRedis{Client: client}
	return &rateLimiter
}

func (rlim *RateLimitRepositoryRedis) SetRequestCounter(session *web_session.WebSession) error {
	ctx := context.Background()
	counterKey, maxRequest := session.GetRequestCounterId(), uint(session.GetRequestsLimitInSeconds())
	err := rlim.Client.Set(ctx, counterKey, maxRequest, 0).Err()
	if err != nil {
		return err
	}
	return rlim.Client.Set(ctx, session.GetRequestTimerId(), time.Now().Unix(), 0).Err()

}

func (rlim *RateLimitRepositoryRedis) GetLastRequestTime(session *web_session.WebSession) (int64, error) {
	ctx := context.Background()
	resetTimeKey := session.GetRequestTimerId()
	lastResetTimeStr, err := rlim.Client.Get(ctx, resetTimeKey).Result()
	if errors.Is(err, redis.Nil) {
		return -1, nil
	}
	lastResetTime, err := strconv.ParseInt(lastResetTimeStr, 10, 64)

	return lastResetTime, err
}

func (rlim *RateLimitRepositoryRedis) DecreaseTokenBucket(session *web_session.WebSession) (bool, error) {
	ctx := context.Background()
	counterKey := session.GetRequestCounterId()
	// Transactional function, optimistic lock.
	txf := func(tx *redis.Tx) error {
		// Get the current value or zero.
		remaingRequests, err := tx.Get(ctx, counterKey).Int()
		if err != nil && err != redis.Nil {
			return err
		}
		if remaingRequests <= 0 {
			throttledError := ThrottledError{}
			return throttledError.ThrottledError()
		}

		// Actual operation (local in optimistic lock).
		remaingRequests--
		// Operation is commited only if the watched keys remain unchanged.
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, counterKey, remaingRequests, 0)
			return nil
		})
		return err
	}
	var err error
	throttledError := ThrottledError{}
	lastResetTime, err := rlim.GetLastRequestTime(session)
	if err != nil {
		return false, err
	}
	if (time.Now().Unix() - lastResetTime) >= session.GetExpireInSeconds() {
		// if elapsed, reset the counter
		err = rlim.SetRequestCounter(session)
		return false, err
	}
	for i := 0; i < int(session.GetRequestsLimitInSeconds()); i++ {

		err = rlim.Client.Watch(ctx, txf, counterKey) // Will return error if not possible.
		if errors.Is(err, redis.Nil) || err == nil {
			return false, nil
		}
		if err != nil && err.Error() == throttledError.ThrottledError().Error() {
			return true, err
		}
		log.Info().Msgf("Optimistic lock failed. %d.", i)
		log.Info().Msg(err.Error())
	}
	return err != nil, err
}
