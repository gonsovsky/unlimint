package storage

import (
	"../shared"
	"github.com/go-redis/redis"
	"log"
	"sync"
	"sync/atomic"
)

type RedisClient struct {
	*redis.Client
	count int32
}

var redisClient *RedisClient
var once sync.Once

const key = "googlehit_"

func NewRedisClient(cfg shared.RedisCfg) *RedisClient {
	once.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     cfg.HostAndPort(),
			Password: cfg.Auth,
			DB:       cfg.Db,
		})

		redisClient = &RedisClient{client,0}
		_, err := redisClient.Ping().Result()
		if err != nil {
			log.Fatalf("Could not connect to redis %v", err)
		}
		client.FlushAll()
	})

	return redisClient
}

func (c *RedisClient) Post(hit shared.GoogleHit) error{
	err := c.Do("SET",key + hit.DocumentPath , hit.ToJSON())
	if err.Val() != "OK" {
		return err.Err()
	}
	atomic.AddInt32(&c.count, 1)
	return nil
}

func (c *RedisClient)  GetCount() int32 {
	return atomic.LoadInt32(&c.count)
}

