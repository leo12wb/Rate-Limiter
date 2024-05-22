package limiter_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/allanmaral/go-expert-rate-limiter-challenge/internal/testutils"
	"github.com/allanmaral/go-expert-rate-limiter-challenge/pkg/limiter"
	"github.com/redis/go-redis/v9"
)

func combine(options ...limiter.Option) []limiter.Option {
	return options
}

func timeAdding(d time.Duration) int64 {
	return time.Now().UTC().Add(d).Truncate(time.Second).Unix()
}

func toInt64(t testing.TB, s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		t.Errorf("failed to convert value %s to int64: %s", s, err)
	}
	return i
}

func TestLimiter(t *testing.T) {
	type response struct {
		httpCode  int
		limit     int
		remaining int
		reset     int64
	}
	type test struct {
		name      string
		config    []limiter.Option
		responses []response
	}
	tests := []test{
		{
			name:   "should not block requests on number of requests bellow the limit",
			config: combine(limiter.WithDefaultLimit(3, 3*time.Second)),
			responses: []response{
				{httpCode: 200, limit: 3, remaining: 2, reset: timeAdding(time.Second)},
				{httpCode: 200, limit: 3, remaining: 1, reset: timeAdding(time.Second)},
			},
		},
		{
			name:   "should not block requests on number of requests equals the limit",
			config: combine(limiter.WithDefaultLimit(4, 3*time.Second)),
			responses: []response{
				{httpCode: 200, limit: 4, remaining: 3, reset: timeAdding(time.Second)},
				{httpCode: 200, limit: 4, remaining: 2, reset: timeAdding(time.Second)},
				{httpCode: 200, limit: 4, remaining: 1, reset: timeAdding(time.Second)},
				{httpCode: 200, limit: 4, remaining: 0, reset: timeAdding(time.Second)},
			},
		},
		{
			name:   "should block requests on number of requests over the limit",
			config: combine(limiter.WithDefaultLimit(2, 3*time.Second)),
			responses: []response{
				{httpCode: 200, limit: 2, remaining: 1, reset: timeAdding(time.Second)},
				{httpCode: 200, limit: 2, remaining: 0, reset: timeAdding(time.Second)},
				{httpCode: 429, limit: 2, remaining: 0, reset: timeAdding(3 * time.Second)},
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			router := limiter.RateLimit(tt.config...)(h)

			for _, response := range tt.responses {
				req := httptest.NewRequest("GET", "/", nil)
				recorder := httptest.NewRecorder()
				router.ServeHTTP(recorder, req)

				result := recorder.Result()
				if response.httpCode != result.StatusCode {
					t.Errorf("expected response.StatusCode(%v) = %v, got %v instead", i, response.httpCode, result.StatusCode)
				}
				if limit := result.Header.Get("X-RateLimit-Limit"); fmt.Sprintf("%d", response.limit) != limit {
					t.Errorf("expected response.X-RateLimit-Limit(%v) = %v, got %v instead", i, response.limit, limit)
				}
				if remaining := result.Header.Get("X-RateLimit-Remaining"); fmt.Sprintf("%d", response.remaining) != remaining {
					t.Errorf("expected response.X-RateLimit-Remaining(%v) = %v, got %v instead", i, response.limit, remaining)
				}
				if reset := result.Header.Get("X-RateLimit-Reset"); toInt64(t, reset) < response.reset {
					t.Errorf("expected response.X-RateLimit-Reset(%v) >= %v, got %v instead", i, response.reset, toInt64(t, reset))
				}
			}
		})
	}
}

func TestLimiterWithRedis(t *testing.T) {
	type response struct {
		httpCode  int
		limit     int
		remaining int
		reset     int64
	}
	type test struct {
		name      string
		config    []limiter.Option
		responses []response
	}
	tests := []test{
		{
			name:   "should not block requests on number of requests bellow the limit",
			config: combine(limiter.WithDefaultLimit(3, 3*time.Second)),
			responses: []response{
				{httpCode: 200, limit: 3, remaining: 2, reset: timeAdding(time.Second)},
				{httpCode: 200, limit: 3, remaining: 1, reset: timeAdding(time.Second)},
			},
		},
		{
			name:   "should not block requests on number of requests equals the limit",
			config: combine(limiter.WithDefaultLimit(4, 3*time.Second)),
			responses: []response{
				{httpCode: 200, limit: 4, remaining: 3, reset: timeAdding(time.Second)},
				{httpCode: 200, limit: 4, remaining: 2, reset: timeAdding(time.Second)},
				{httpCode: 200, limit: 4, remaining: 1, reset: timeAdding(time.Second)},
				{httpCode: 200, limit: 4, remaining: 0, reset: timeAdding(time.Second)},
			},
		},
		{
			name:   "should block requests on number of requests over the limit",
			config: combine(limiter.WithDefaultLimit(2, 3*time.Second)),
			responses: []response{
				{httpCode: 200, limit: 2, remaining: 1, reset: timeAdding(time.Second)},
				{httpCode: 200, limit: 2, remaining: 0, reset: timeAdding(time.Second)},
				{httpCode: 429, limit: 2, remaining: 0, reset: timeAdding(3 * time.Second)},
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redisContainer := testutils.CreateRedisTestingContainer(t)
			redisClient := redis.NewClient(&redis.Options{
				Addr: redisContainer.Endpoint,
			})
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			config := append(tt.config, limiter.WithRedisStore(redisClient))
			router := limiter.RateLimit(config...)(h)

			for _, response := range tt.responses {
				req := httptest.NewRequest("GET", "/", nil)
				recorder := httptest.NewRecorder()
				router.ServeHTTP(recorder, req)

				result := recorder.Result()
				if response.httpCode != result.StatusCode {
					t.Errorf("expected response.StatusCode(%v) = %v, got %v instead", i, response.httpCode, result.StatusCode)
				}
				if limit := result.Header.Get("X-RateLimit-Limit"); fmt.Sprintf("%d", response.limit) != limit {
					t.Errorf("expected response.X-RateLimit-Limit(%v) = %v, got %v instead", i, response.limit, limit)
				}
				if remaining := result.Header.Get("X-RateLimit-Remaining"); fmt.Sprintf("%d", response.remaining) != remaining {
					t.Errorf("expected response.X-RateLimit-Remaining(%v) = %v, got %v instead", i, response.limit, remaining)
				}
				if reset := result.Header.Get("X-RateLimit-Reset"); toInt64(t, reset) < response.reset {
					t.Errorf("expected response.X-RateLimit-Reset(%v) >= %v, got %v instead", i, response.reset, toInt64(t, reset))
				}
			}
		})
	}
}
