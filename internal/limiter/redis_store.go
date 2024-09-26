package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RateLimiter struct
type RateLimiter struct {
	ipStore    Limiter
	tokenStore Limiter
	ipLimit    int
	tokenLimit int
	blockTime  time.Duration
}

// Limiter interface
type Limiter interface {
	AllowRequest(key string, limit int, blockTime time.Duration) (bool, error)
}

// RedisStore struct
type RedisStore struct {
	client *redis.Client
}

// NewRedisStore creates a new RedisStore

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client}
}

// AllowRequest checks if a request is allowed
func (rs *RedisStore) AllowRequest(key string, limit int, blockTime time.Duration) (bool, error) {
	ctx := context.Background()
	val, err := rs.client.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return false, err
	}

	fmt.Printf("Current value for key %s: %d\n", key, val)

	if val >= limit {
		return false, nil
	}

	pipe := rs.client.TxPipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, blockTime)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	fmt.Printf("Updated value for key %s: %d\n", key, val+1)
	return true, nil
}

// NewRateLimiter creates a new RateLimiter
func NewRateLimiter(ipStore, tokenStore Limiter, ipLimit, tokenLimit int, blockTime time.Duration) *RateLimiter {
	return &RateLimiter{
		ipStore:    ipStore,
		tokenStore: tokenStore,
		ipLimit:    ipLimit,
		tokenLimit: tokenLimit,
		blockTime:  blockTime,
	}
}

// GetLimiter returns the appropriate limiter instance, key, and limit based on IP or token
func (rl *RateLimiter) GetLimiter(ip, token string) (Limiter, string, int) {
	if token != "" {
		return rl.tokenStore, token, rl.tokenLimit
	}
	return rl.ipStore, ip, rl.ipLimit
}

// GetBlockTime returns the block time duration
func (rl *RateLimiter) GetBlockTime() time.Duration {
	return rl.blockTime
}
