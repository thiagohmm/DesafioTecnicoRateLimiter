package limiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client}
}

func (rs *RedisStore) Allow(key string, limit int, blockTime time.Duration) (bool, error) {
	ctx := context.Background()
	count, err := rs.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		rs.client.Expire(ctx, key, time.Second)
	}

	if int(count) > limit {
		rs.client.Set(ctx, key+":blocked", true, blockTime)
		return false, nil
	}

	return true, nil
}
