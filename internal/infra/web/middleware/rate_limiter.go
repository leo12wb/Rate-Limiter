package middleware

import (
	"context"
	"github.com/aluferraz/go-expert-rate-limiter/internal/entity/web_session"
	"github.com/aluferraz/go-expert-rate-limiter/internal/infra/persistence/rate_limit"
	"github.com/aluferraz/go-expert-rate-limiter/internal/value_objects"
	"net/http"
)

type RateLimiterMiddleware struct {
	Storage       rate_limit.RateLimitRepository
	IPLimit       uint
	ApiKeyLimit   uint
	ExpireSeconds uint
}

func NewRateLimiter(storage *rate_limit.RateLimitRepository, requestLimits value_objects.RequestLimits) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		*storage,
		requestLimits.IPLimit,
		requestLimits.APILimit,
		requestLimits.ExpireSeconds,
	}
}

// HTTP middleware setting a value on the request context
func (rl *RateLimiterMiddleware) RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		websession, _ := web_session.NewWebSession(
			r.RemoteAddr,
			r.Header.Get("API_KEY"),
			value_objects.NewRequestLimit(rl.IPLimit, rl.ApiKeyLimit, rl.ExpireSeconds),
		)

		// create new context from `r` request context, and assign key `"user"`
		// to value of `"123"`
		ctx := context.WithValue(r.Context(), "websession", websession)

		denied, err := rl.Storage.DecreaseTokenBucket(&websession)

		if denied {
			http.Error(w, err.Error(), 429)
			return

		} else {
			if err != nil {
				w.WriteHeader(500)
			}
		}
		// call the next handler in the chain, passing the response writer and
		// the updated request object with the new context value.
		//
		// note: context.Context values are nested, so any previously set
		// values will be accessible as well, and the new `"user"` key
		// will be accessible from this point forward.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
