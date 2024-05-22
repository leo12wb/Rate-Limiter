package limiter

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type Option func(*Config)

type Config struct {
	Limits       map[string]Limit
	DefaultLimit Limit
	Store        CounterStore
}

type Limit struct {
	RequestsPerSecond int64
	TimeBlocked       time.Duration
}

func newConfig(options ...Option) Config {
	config := Config{
		Limits: make(map[string]Limit),
		Store:  newLocalStore(),
	}
	for _, opt := range options {
		opt(&config)
	}

	return config
}

func WithDefaultLimit(requestPerSecond int64, timeBlocked time.Duration) Option {
	return func(c *Config) {
		c.DefaultLimit.RequestsPerSecond = requestPerSecond
		c.DefaultLimit.TimeBlocked = timeBlocked
	}
}

func WithIPLimits(limits map[string]Limit) Option {
	return func(c *Config) {
		for k, v := range limits {
			c.Limits[IpKey(k)] = v
		}
	}
}

func WithTokenLimits(limits map[string]Limit) Option {
	return func(c *Config) {
		for k, v := range limits {
			c.Limits[TokenKey(k)] = v
		}
	}
}

func WithRedisStore(client *redis.Client) Option {
	return func(c *Config) {
		c.Store = NewRedisStore(client)
	}
}
