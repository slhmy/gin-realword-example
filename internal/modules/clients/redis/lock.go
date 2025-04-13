package redis_client

import (
	"sync"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
)

// See more: https://github.com/go-redsync/redsync

var (
	lockInitMutex sync.Mutex
	redsyncClient *redsync.Redsync
)

func GetLockClient() *redsync.Redsync {
	if redsyncClient == nil {
		lockInitMutex.Lock()
		defer lockInitMutex.Unlock()
		if redsyncClient != nil {
			return redsyncClient
		}

		pool := goredis.NewPool(GetRDB())
		redsyncClient = redsync.New(pool)
	}
	return redsyncClient
}
