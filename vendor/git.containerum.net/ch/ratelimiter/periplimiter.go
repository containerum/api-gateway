package ratelimiter

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

var incrExpire = redis.NewScript(`
	local current
	current = redis.call("incr",KEYS[1])
	if tonumber(current) == 1 then
    	redis.call("expire",KEYS[1],1)
	end
`)

type PerIPLimiterConfig struct {
	RedisAddress  string
	RedisPassword string
	RedisDB       int
	Timeout       time.Duration
	RateLimit     int
}

// PerIPLimiter limits request rate for single IP
type PerIPLimiter struct {
	redisClient *redis.Client
	rateLimit   uint64
}

func NewPerIPLimiter(cfg *PerIPLimiterConfig) (*PerIPLimiter, error) {
	if cfg.RateLimit <= 0 {
		return nil, errors.New("invalid rate limit")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddress,
		Password:     cfg.RedisPassword,
		DB:           cfg.RedisDB,
		DialTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		ReadTimeout:  cfg.Timeout,
	})
	if err := redisClient.Ping().Err(); err != nil {
		return nil, err
	}
	return &PerIPLimiter{
		redisClient: redisClient,
		rateLimit:   uint64(cfg.RateLimit),
	}, nil
}

func (l *PerIPLimiter) Limit(ip string) (allowCall bool, err error) {
	requests := l.redisClient.Get(ip)
	n, err := requests.Uint64()
	if err != nil && err != redis.Nil {
		return true, fmt.Errorf("value to int64 cast: %v", err)
	}
	if n > l.rateLimit {
		return false, nil
	}
	if err := incrExpire.Run(l.redisClient, []string{ip}).Err(); err != nil && err != redis.Nil {
		return true, fmt.Errorf("increment with expire run: %v", err)
	}
	return true, nil
}
