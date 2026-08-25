package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
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

func (r *UserRedisRepo) Get(ctx context.Context, userID uuid.UUID) (*entities.User, error) {
	val, err := r.client.Get(ctx, "user:"+userID.String()).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var user entities.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, err
	}
	return &user, nil

}

func (r *UserRedisRepo) Set(ctx context.Context, req *entities.User) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, "user:"+req.ID.String(), data, 10*time.Minute).Err()
}
