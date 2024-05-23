package usecase

import (
	"context"
	"os"
	"time"

	"github.com/leo12wb/Rate-Limiter/ratelimiter/internal/entity"
	"github.com/leo12wb/Rate-Limiter/ratelimiter/internal/infra/inmemory"
	"github.com/leo12wb/Rate-Limiter/ratelimiter/internal/infra/rdb"

	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var (
	defaultConfig = entity.Config{}
)

func ConfigLimiter(database entity.DatabaseRepository) {
	ctx := context.Background()
	for _, token := range defaultConfig.Limiter.Tokens {
		rateLimiterInfo := entity.RateLimiterInfo{
			Key:      token.Token,
			Requests: token.Requests,
			Every:    token.Every,
		}

		if err := database.Create(ctx, rateLimiterInfo, 0); err != nil {
			log.Println(err)
		}
	}

	for _, ip := range defaultConfig.Limiter.IPS {
		rateLimiterInfo := entity.RateLimiterInfo{
			Key:      ip.IP,
			Requests: ip.Requests,
			Every:    ip.Every,
		}

		if err := database.Create(ctx, rateLimiterInfo, 0); err != nil {
			log.Println(err)
		}
	}
}

func LoadConfig() {
	if fileExists(".env") {
		loadConfigFromEnv()
		return
	}
	loadConfigFromJson()
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func loadConfigFromJson() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&defaultConfig)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct: %w", err))
	}
}

func loadConfigFromEnv() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.MergeInConfig(); err != nil {
		panic(fmt.Errorf("Erro ao ler o arquivo .env: %s", err))
	}

	defaultConfig.Limiter.Database.InMemory = viper.GetBool("LIMITER_DATABASE_INMEMORY")
	defaultConfig.Limiter.Database.Redis = viper.GetBool("LIMITER_DATABASE_REDIS")
	defaultConfig.Limiter.Default.Requests = viper.GetInt("LIMITER_DEFAULT_REQUESTS")
	defaultConfig.Limiter.Default.Every = viper.GetInt("LIMITER_DEFAULT_EVERY")

	for i := 0; ; i++ {
		ipKey := fmt.Sprintf("LIMITER_IPS_%d_IP", i)
		if !viper.IsSet(ipKey) {
			break
		}

		ip := viper.GetString(ipKey)
		requests := viper.GetInt(fmt.Sprintf("LIMITER_IPS_%d_REQUESTS", i))
		every := viper.GetInt(fmt.Sprintf("LIMITER_IPS_%d_EVERY", i))

		defaultConfig.Limiter.IPS = append(defaultConfig.Limiter.IPS, entity.IP{
			IP:       ip,
			Requests: requests,
			Every:    every,
		})
	}

	for i := 0; ; i++ {
		tokenKey := fmt.Sprintf("LIMITER_TOKENS_%d_TOKEN", i)
		if !viper.IsSet(tokenKey) {
			break
		}

		token := viper.GetString(tokenKey)
		requests := viper.GetInt(fmt.Sprintf("LIMITER_TOKENS_%d_REQUESTS", i))
		every := viper.GetInt(fmt.Sprintf("LIMITER_TOKENS_%d_EVERY", i))

		defaultConfig.Limiter.Tokens = append(defaultConfig.Limiter.Tokens, entity.Token{
			Token:    token,
			Requests: requests,
			Every:    every,
		})
	}
}

func GetStorage() entity.DatabaseRepository {
	if defaultConfig.Limiter.Database.InMemory {
		return inmemory.NewDatabaseRepository()
	}

	if defaultConfig.Limiter.Database.Redis {
		database := redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "",
			DB:       0,
		})
		return rdb.NewDatabaseRepository(database)
	}

	panic("storage must be provide (inMemory/redis or implement yourself)")
}

func CheckLimit(ctx context.Context, database entity.DatabaseRepository, key string) bool {
	config, err := database.Read(ctx, key)

	if err != nil {
		database.Create(ctx, entity.RateLimiterInfo{
			Key:      key,
			Requests: defaultConfig.Limiter.Default.Requests,
			Every:    defaultConfig.Limiter.Default.Every,
		}, 0)

		return CheckLimit(ctx, database, key)
	}

	limiter, err := database.CheckLimit(ctx, key, config.Requests, time.Duration(config.Every)*time.Second)

	if err != nil {
		log.Println(err)
		return false
	}

	return limiter
}
