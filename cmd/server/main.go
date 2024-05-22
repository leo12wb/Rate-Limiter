package main

import (
	"fmt"
	"net/http"

	"github.com/allanmaral/go-expert-rate-limiter-challenge/configs"
	"github.com/allanmaral/go-expert-rate-limiter-challenge/pkg/limiter"
	"github.com/redis/go-redis/v9"
)

func main() {
	config := configs.Load()
	rc := redis.NewClient(&redis.Options{Addr: config.RedisAddr})

	limit := limiter.RateLimit(
		limiter.WithRedisStore(rc),
		limiter.WithDefaultLimit(config.DefaultLimit.RequestsPerSecond, config.DefaultLimit.TimeBlocked),
		limiter.WithIPLimits(config.IPLimits),
		limiter.WithTokenLimits(config.TokenLimits),
	)

	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	fmt.Printf("Listening on http://localhost:%s\n", config.Port)
	http.ListenAndServe(fmt.Sprintf(":%s", config.Port), limit(router))
}
