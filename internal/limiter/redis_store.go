package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

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
