package limiter

import (
	"time"
)

type Limiter interface {
	AllowRequest(ip string, token string) (bool, error)
}

type RateLimiter struct {
	ipLimiter    LimiterStore
	tokenLimiter LimiterStore
	ipLimit      int
	tokenLimit   int
	blockTime    time.Duration
}

func NewRateLimiter(ipLimiter, tokenLimiter LimiterStore, ipLimit, tokenLimit int, blockTime time.Duration) *RateLimiter {
	return &RateLimiter{
		ipLimiter:    ipLimiter,
		tokenLimiter: tokenLimiter,
		ipLimit:      ipLimit,
		tokenLimit:   tokenLimit,
		blockTime:    blockTime,
	}
}

func (rl *RateLimiter) AllowRequest(ip string, token string) (bool, error) {
	if token != "" {
		return rl.tokenLimiter.Allow(token, rl.tokenLimit, rl.blockTime)
	}
	return rl.ipLimiter.Allow(ip, rl.ipLimit, rl.blockTime)
}

type LimiterStore interface {
	Allow(key string, limit int, blockTime time.Duration) (bool, error)
}
