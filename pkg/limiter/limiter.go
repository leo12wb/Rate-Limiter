package limiter

import (
	"fmt"
	"net"
	"net/http"
)

type RateLimiter struct {
	counter *limitCounter
}

func NewRateLimiter(options ...Option) *RateLimiter {
	return &RateLimiter{counter: newLimitCounter(options...)}
}

func RateLimit(options ...Option) func(next http.Handler) http.Handler {
	return NewRateLimiter(options...).Handler
}

func (rl *RateLimiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("API_KEY")
		ip := requestIP(r)
		info, err := rl.counter.Increment(token, ip)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", info.Limit))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", info.Remaining))
		w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", info.Reset.Unix()))
		if info.Remaining < 0 {
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("RetryAfter", fmt.Sprintf("%d", info.Reset.Unix()))
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func requestIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}
	return ip
}
