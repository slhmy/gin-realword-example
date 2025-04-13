package redis_client

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/go-redis/cache/v9"
)

const (
	cacheRefreshEventChannel = "cacheRefreshEventChannel"
)

var (
	cacheInitMutex sync.Mutex
	cacheClient    *cache.Cache
)

func getCache() *cache.Cache {
	if cacheClient == nil {
		cacheInitMutex.Lock()
		defer cacheInitMutex.Unlock()
		if cacheClient != nil {
			return cacheClient
		}

		cacheClient = cache.New(&cache.Options{
			Redis:      GetRDB(),
			LocalCache: cache.NewTinyLFU(1000, time.Minute),
		})
	}
	return cacheClient
}

func SubscribeCacheRefreshEvent(ctx context.Context) {
	_, err := GetRDB().Set(ctx, cacheRefreshEventChannel, cacheRefreshEventChannel, 0).Result()
	if err != nil {
		panic(err)
	}

	go func() {
		pubsub := GetRDB().Subscribe(ctx, cacheRefreshEventChannel)
		defer func() {
			err := pubsub.Close()
			if err != nil {
				slog.Error("Error closing pubsub", "error", err)
			}
		}()

		ch := pubsub.Channel()
		for {
			select {
			case msg := <-ch:
				if msg.Payload == cacheRefreshEventChannel {
					slog.Info("Cache refresh event received", "key", msg.Payload)
					getCache().DeleteFromLocalCache(msg.Payload)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func PublishCacheRefreshEvent(ctx context.Context, key string) error {
	return GetRDB().Publish(ctx, cacheRefreshEventChannel, key).Err()
}

func GetCache[T any](ctx context.Context, key string) (*T, error) {
	var value T
	err := getCache().Get(ctx, key, &value)
	return &value, err
}

func SetCache(ctx context.Context, key string, value any, expiration time.Duration) error {
	if err := getCache().Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   expiration,
	}); err != nil {
		return err
	}
	return PublishCacheRefreshEvent(ctx, key)
}

func DelCache(ctx context.Context, key string) error {
	if err := getCache().Delete(ctx, key); err != nil {
		return err
	}
	return PublishCacheRefreshEvent(ctx, key)
}
