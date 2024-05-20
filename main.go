package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

// RateLimiterMiddleware é um middleware HTTP para rate limiting
type RateLimiterMiddleware struct {
	maxRequestsPerSecond int           // Número máximo de requisições por segundo
	blockTime            time.Duration // Tempo de bloqueio em caso de exceder o limite (em segundos)
	redisClient          *redis.Client
}

// NewRateLimiterMiddleware cria um novo middleware rate limiter
func NewRateLimiterMiddleware(maxRequestsPerSecond int, blockTime time.Duration, redisAddr, redisPassword string, redisDB int) *RateLimiterMiddleware {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})
	return &RateLimiterMiddleware{
		maxRequestsPerSecond: maxRequestsPerSecond,
		blockTime:            blockTime,
		redisClient:          client,
	}
}

// Middleware retorna um middleware HTTP que limita o tráfego com base no endereço IP
func (rlm *RateLimiterMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		// Verifica se o IP excede o limite
		if rlm.exceedsLimit(ip) {
			http.Error(w, "you have reached the maximum number of requests allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}

		// Executa o próximo handler
		next.ServeHTTP(w, r)
	})
}

// exceedsLimit verifica se o IP excede o limite de requisições
func (rlm *RateLimiterMiddleware) exceedsLimit(ip string) bool {
	// Verifica se o IP está bloqueado
	blocked, _ := rlm.redisClient.Exists(r.Context(), "blocked:"+ip).Result()
	if blocked == 1 {
		return true
	}

	// Verifica o número de requisições feitas pelo IP
	count, _ := rlm.redisClient.Incr(r.Context(), "requests:"+ip).Result()
	if count == 1 {
		// Define a expiração da chave no Redis
		rlm.redisClient.Expire(r.Context(), "requests:"+ip, rlm.blockTime)
	}

	// Verifica se o número de requisições excede o limite
	if count > int64(rlm.maxRequestsPerSecond) {
		// Bloqueia o IP por um determinado tempo
		rlm.redisClient.Set(r.Context(), "blocked:"+ip, true, rlm.blockTime)
		return true
	}

	return false
}

func main() {
	// Carrega as variáveis de ambiente do arquivo .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Erro ao carregar o arquivo .env")
		os.Exit(1)
	}

	// Configurações do rate limiter
	maxRequestsPerSecond, _ := strconv.Atoi(os.Getenv("MAX_REQUESTS_PER_SECOND"))
	blockTimeInSeconds, _ := strconv.Atoi(os.Getenv("BLOCK_TIME_SECONDS"))
	blockTime := time.Duration(blockTimeInSeconds) * time.Second
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	// Criar um novo middleware rate limiter
	limiter := NewRateLimiterMiddleware(maxRequestsPerSecond, blockTime, redisAddr, redisPassword, redisDB)

	// Configurar rota com o middleware do rate limiter
	http.Handle("/", limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})))

	// Iniciar o servidor
	fmt.Println("Servidor iniciado na porta 8080...")
	http.ListenAndServe(":8080", nil)
}
