package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type UserCacheRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewUserCacheRepository(client *redis.Client, ctx context.Context) *UserCacheRepository {
	return &UserCacheRepository{
		client: client,
		ctx:    context.Background(),
	}
}
