package util

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

type RedSync struct {
	mutex *redsync.Mutex
}

// NewRedSync return a *RedSync
func NewRedSync(client *redis.Client, mutexName string) *RedSync {
	// Create a pool with go-redis (or redigo) which is the pool redisync will
	// use while communicating with Redis. This can also be any pool that
	// implements the `redis.Pool` interface.
	pool := goredis.NewPool(client)
	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs := redsync.New(pool)
	return &RedSync{mutex: rs.NewMutex(mutexName)}
}

// Lock return a bool, error
func (r *RedSync) Lock(ctx context.Context) (bool, error) {
	// Obtain a lock for our given mutex. After this is successful, no one else
	// can obtain the same lock (the same mutex name) until we unlock it.
	if err := r.mutex.LockContext(ctx); err != nil {
		return false, err
	}
	return true, nil
}

// UnLock return a bool,error
func (r *RedSync) UnLock(ctx context.Context) (bool, error) {
	// Do your work that requires the lock.

	// Release the lock so other processes or threads can obtain a lock.
	if ok, err := r.mutex.UnlockContext(ctx); !ok || err != nil {
		return false, err
	}
	return true, nil
}
