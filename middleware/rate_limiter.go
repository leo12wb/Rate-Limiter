package middleware

import (
	"go-rate-limiter/limiter"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RateLimiter() echo.MiddlewareFunc {
	rateLimiter, err := limiter.NewRateLimiter()
	if err != nil {
		panic(err)
	}

	config := middleware.RateLimiterConfig{
		Skipper:      middleware.DefaultSkipper,
		Store:        rateLimiter,
		IdentifierExtractor: func(c echo.Context) (string, error) {
			token := c.Request().Header.Get("API_KEY")
			if token != "" {
				return token, nil
			}
			return c.RealIP(), nil
		},
		LimitReachedHandler: func(c echo.Context) error {
			return c.JSON(429, map[string]string{"message": "you have reached the maximum number of requests or actions allowed within a certain time frame"})
		},
	}
	return middleware.RateLimiterWithConfig(config)
}
