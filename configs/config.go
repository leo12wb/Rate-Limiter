package configs

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/allanmaral/go-expert-rate-limiter-challenge/pkg/limiter"
	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	RedisAddr    string
	DefaultLimit limiter.Limit
	IPLimits     map[string]limiter.Limit
	TokenLimits  map[string]limiter.Limit
}

func Load() Config {
	_ = godotenv.Load()

	requestsPerSecond := getInt64Env("REQ_PER_SECOND")
	timeBlocked := getInt64Env("TIME_BLOCKED")
	defaultLimit := limiter.Limit{
		RequestsPerSecond: requestsPerSecond,
		TimeBlocked:       time.Duration(timeBlocked) * time.Second,
	}

	ipLimits := getLimitEnv("IP_LIMITS")
	tokenLimits := getLimitEnv("TOKEN_LIMITS")

	config := Config{
		Port:         os.Getenv("PORT"),
		RedisAddr:    os.Getenv("REDIS_ADDR"),
		DefaultLimit: defaultLimit,
		IPLimits:     ipLimits,
		TokenLimits:  tokenLimits,
	}

	if config.Port == "" {
		config.Port = "8080"
	}

	return config
}

func getInt64Env(name string) int64 {
	value := os.Getenv(name)
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Fatalf("invalid %s '%q': %s", name, value, err)
	}
	return i
}

func getLimitEnv(name string) map[string]limiter.Limit {
	value := os.Getenv(name)
	limits := make(map[string]limiter.Limit)

	if value == "" {
		return limits
	}

	entries := strings.Split(value, ";")
	for _, e := range entries {
		parts := strings.Split(e, ",")
		key := parts[0]
		rps, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			log.Fatalf("invalid limit %s, could not parse RequestPerSecond for Key '%q': %s", name, key, err)
		}
		tb, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			log.Fatalf("invalid limit %s, could not parse TimeBlocked for Key '%q': %s", name, key, err)
		}
		limits[key] = limiter.Limit{RequestsPerSecond: rps, TimeBlocked: time.Duration(tb) * time.Second}
	}
	return limits
}
