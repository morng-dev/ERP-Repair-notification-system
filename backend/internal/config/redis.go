package config

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func SetupRedis(cfg *Config) *redis.Client {
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	log.Println("Redis connecting:", addr)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal("Redis error:", err)
	}

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}

	return client
}
