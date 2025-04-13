package redis_client

import (
	"sync"

	"github.com/go-redis/redis_rate/v10"
)

// See more: https://github.com/go-redis/redis_rate

var (
	limiterInitMutex sync.Mutex
	limiter          *redis_rate.Limiter
)

func GetRateLimitClient() *redis_rate.Limiter {
	if limiter == nil {
		limiterInitMutex.Lock()
		defer limiterInitMutex.Unlock()
		if limiter != nil {
			return limiter
		}

		limiter = redis_rate.NewLimiter(GetRDB())
	}
	return limiter
}
