package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func SetupRedis(cfg *Config) *redis.Client {
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	log.Println("Redis connecting:", addr)
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     cfg.RedisPass,
		DB:           cfg.RedisDB,
		PoolSize:     20,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		MaxRetries:   3,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Redis unavailable: %v", err)
		log.Println("Application will continue without Redis.")
	} else {
		log.Println("Redis connected successfully")
	}

	return client
}
