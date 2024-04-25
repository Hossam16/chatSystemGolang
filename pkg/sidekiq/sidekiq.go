package sidekiq

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/bsm/redislock"
	"github.com/Hossam16/go-chat-creation-api/configs"
)
var redisClient *redis.Client
var redisLocker *redislock.Client

type sidekiqJob struct {
	Class string   `json:"class"`
	Args  []string `json:"args"`
	Retry bool     `json:"retry"`
	Queue string   `json:"queue"`
}

func Push(queue string, class string, args ...string) error {
	job := sidekiqJob {
		Class: class,
		Args:  args,
		Queue: queue,
		Retry: true,
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     configs.RedisAddress,
		Password: "",
		DB:       0,
	})
	redisLocker = redislock.New(redisClient)

	// redisClient, err := redis.GetRedis()
	// if err != nil {
	// 	return err
	// }

	jobBytes, err := json.Marshal(job)
	if err != nil {
		return err
	}
	ctx := context.Background()
	_, err = redisClient.ZAdd(ctx,"schedule",redis.Z{
		Member: jobBytes,
	}).Result()
	return err
}
