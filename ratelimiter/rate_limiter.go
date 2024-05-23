package ratelimiter

import (
	"net/http"
	"strings"

	"github.com/leo12wb/Rate-Limiter/ratelimiter/internal/entity"
	"github.com/leo12wb/Rate-Limiter/ratelimiter/internal/usecase"
)

var limiterOpts = &LimiterOpts{}

type LimiterOpts struct {
	Storage entity.DatabaseRepository
}

type LimiterOption func(d *LimiterOpts)

func Storage(value entity.DatabaseRepository) LimiterOption {
	return func(c *LimiterOpts) {
		c.Storage = value
	}
}

func Initialize(opts ...LimiterOption) {
	usecase.LoadConfig()

	for _, opt := range opts {
		opt(limiterOpts)
	}

	if limiterOpts.Storage == nil {
		limiterOpts.Storage = usecase.GetStorage()
	}

	usecase.ConfigLimiter(limiterOpts.Storage)
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := getIP(r)
		queryValues := r.URL.Query()
		if token, ok := queryValues["token"]; ok {
			key = token[0]
		}

		limiter := usecase.CheckLimit(r.Context(), limiterOpts.Storage, key)

		if !limiter {
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	ip := strings.Split(r.RemoteAddr, ":")[0]
	return ip
}
