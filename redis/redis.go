package redisScrapper

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func ConnectToRedis() *redis.Client {

	fmt.Printf("connection to redis %s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))

	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

}
