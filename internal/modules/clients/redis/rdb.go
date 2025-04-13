package redis_client

import (
	"sync"

	"gin-realword-example/internal/modules/core"

	"github.com/redis/go-redis/v9"
)

const (
	redisUrlsConfigKey     = "redis.urls"
	redisPasswordConfigKey = "redis.password"
)

var (
	rdbInitMutex sync.Mutex
	redisUrls    []string
	password     string
)

func init() {
	redisUrls = core.ConfigStore.GetStringSlice(redisUrlsConfigKey)
	password = core.ConfigStore.GetString(redisPasswordConfigKey)
}

var redisClient redis.UniversalClient

func GetRDB() redis.UniversalClient {
	if redisClient == nil {
		rdbInitMutex.Lock()
		defer rdbInitMutex.Unlock()
		if redisClient != nil {
			return redisClient
		}

		if len(redisUrls) == 0 {
			panic("No redis hosts configured")
		}
		if len(redisUrls) == 1 {
			redisClient = redis.NewClient(&redis.Options{
				Addr:     redisUrls[0],
				Password: password,
			})
		}
		if len(redisUrls) > 1 {
			redisClient = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:    redisUrls,
				Password: password,
			})
		}
	}
	return redisClient
}
