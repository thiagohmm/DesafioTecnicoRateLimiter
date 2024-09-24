package main

import (
	"os"
	"strconv"
	"time"

	"github.com/thiagohmm/DesafioTecnicoRateLimiter/internal/limiter"
	"github.com/thiagohmm/DesafioTecnicoRateLimiter/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	ipLimit, _ := strconv.Atoi(os.Getenv("LIMITER_IP"))
	tokenLimit, _ := strconv.Atoi(os.Getenv("LIMITER_TOKEN"))
	blockDuration, _ := time.ParseDuration(os.Getenv("BLOCK_DURATION"))

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
