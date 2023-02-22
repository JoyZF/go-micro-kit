package model

import (
	"fmt"
	"github.com/JoyZF/go-micro-kit/service/greeter/internal/conf"
	"github.com/redis/go-redis/v9"
	"go-micro.dev/v4/util/log"
)

var (
	rdb *redis.Client
)

func InitRedis(c *conf.Redis) {
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: c.Password, // no password set
		DB:       c.DB,       // use default DB
	})
}

// GetRdb returns a Redis client representing a pool of zero or more underlying connections.
func GetRdb() *redis.Client {
	if rdb == nil {
		log.Fatal("redis init fail %+v", err)
	}
	return rdb
}
