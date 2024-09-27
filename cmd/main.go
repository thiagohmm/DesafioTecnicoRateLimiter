package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/thiagohmm/DesafioTecnicoRateLimiter/internal/limiter"
	"github.com/thiagohmm/DesafioTecnicoRateLimiter/internal/middleware"
)

func main() {
	// Carregar variáveis de ambiente do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file: %v", err)
	}
	fmt.Println(".env file loaded successfully")

	ipLimitStr := os.Getenv("LIMITER_IP")
	if ipLimitStr == "" {
		log.Fatalf("Environment variable LIMITER_IP is not set")
	}
	ipLimit, err := strconv.Atoi(ipLimitStr)
	if err != nil {
		log.Fatalf("Invalid IP limit: %v", err)
	}
	fmt.Printf("IP Limit: %d\n", ipLimit)

	tokenLimitStr := os.Getenv("LIMITER_TOKEN")
	if tokenLimitStr == "" {
		log.Fatalf("Environment variable LIMITER_TOKEN is not set")
	}
	tokenLimit, err := strconv.Atoi(tokenLimitStr)
	if err != nil {
		log.Fatalf("Invalid Token limit: %v", err)
	}
	fmt.Printf("Token Limit: %d\n", tokenLimit)

	blockDurationStr := os.Getenv("BLOCK_DURATION")
	if blockDurationStr == "" {
		log.Fatalf("Environment variable BLOCK_DURATION is not set")
	}
	blockDuration, err := time.ParseDuration(blockDurationStr)
	if err != nil {
		log.Fatalf("Invalid Block Duration: %v", err)
	}
	fmt.Printf("Block Duration: %s\n", blockDuration)

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatalf("Environment variable REDIS_ADDR is not set")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	limiter := limiter.NewRateLimiter(
		limiter.NewRedisStore(rdb),
		limiter.NewRedisStore(rdb),
		ipLimit,
		tokenLimit,
		blockDuration,
	)
	r := gin.Default()
	r.Use(middleware.RateLimiterMiddleware(limiter))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Request successful"})
	})

	r.Run(":8080")
}
