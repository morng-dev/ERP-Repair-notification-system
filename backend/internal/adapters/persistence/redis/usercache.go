package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/cache"
	"github.com/redis/go-redis/v9"
)

type UserRedisRepo struct {
	client *redis.Client
}

func NewUserRedisRepo(client *redis.Client) cache.UserRedisRepo {
	return &UserRedisRepo{client: client}
}

func userListKey(page, limit int) string {
	return fmt.Sprintf("users:list:page:%d:limit:%d", page, limit)
}

func (r *UserRedisRepo) GetAll(ctx context.Context, page, limit int) (*entities.UserListCacheData, error) {
	key := userListKey(page, limit)
	byteData, err := r.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var data entities.UserListCacheData
	if err := json.Unmarshal(byteData, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *UserRedisRepo) SetAll(ctx context.Context, page, limit int, data *entities.UserListCacheData) error {
	key := userListKey(page, limit)
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, value, 5*time.Minute).Err()
}

func (r *UserRedisRepo) DeleteAll(ctx context.Context) error {
	var cursor uint64

	for {
		keys, nextCursor, err := r.client.Scan(ctx, cursor, "users:list:*", 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := r.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}
