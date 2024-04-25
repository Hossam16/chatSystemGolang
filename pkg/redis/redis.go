package redis

import (
	"log"
	"context"
	"github.com/go-redis/redis"
	"github.com/bsm/redislock"
	"github.com/Hossam16/go-chat-creation-api/configs"
)

var redisClient *redis.Client
var redisLocker *redislock.Client

func GetRedis() (*redis.Client, error) {
	if redisClient == nil {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     configs.RedisAddress,
			Password: "",
			DB:       0,
		})
		ctx := context.Background()
		err := redisClient.Ping(ctx).Err()
		if err != nil {
			return nil, err
		}
		redisLocker = redislock.New(redisClient)
	}
	return redisClient, nil
}

func GetLocker() (*redislock.Client) {
	if redisClient == nil {
		log.Fatalln("Redis client is not initialized yet")
	}
	return redisLocker
}