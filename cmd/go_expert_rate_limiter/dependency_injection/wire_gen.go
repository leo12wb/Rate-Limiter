// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package dependency_injection

import (
	"github.com/leo12wb/Rate-Limiter/internal/infra/persistence/rate_limit"
	"github.com/leo12wb/Rate-Limiter/internal/infra/web/middleware"
	"github.com/leo12wb/Rate-Limiter/internal/value_objects"
	"github.com/redis/go-redis/v9"
)

// Injectors from wire.go:

func NewRateLimitMiddleware(client *redis.Client, requestLimits value_objects.RequestLimits) *middleware.RateLimiterMiddleware {
	rateLimitRepository := rate_limit.NewRateLimitRepositoryRedis(client)
	rateLimiterMiddleware := middleware.NewRateLimiter(rateLimitRepository, requestLimits)
	return rateLimiterMiddleware
}
