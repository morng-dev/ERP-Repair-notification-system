package redis

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/cache"
	"github.com/redis/go-redis/v9"
)

type Permissioncache struct {
	client *redis.Client
}

func NewPermissioncache(client *redis.Client) cache.PermissionRedisRepo {
	return &Permissioncache{client: client}
}

func (r *Permissioncache) Set(ctx context.Context, userID uuid.UUID, permissions []*entities.Permission) error {
	key := "user:permissions:" + userID.String()
	data, err := json.Marshal(permissions)
	if err != nil {
		log.Printf("Redis SET ERROR: %v", err)
		return err
	}

	log.Printf("Redis SET SUCCESS: key=%s", key)
	return r.client.Set(ctx, key, data, 5*time.Minute).Err()
}

func (r *Permissioncache) Get(ctx context.Context, userID uuid.UUID) ([]*entities.Permission, error) {
	key := "user:permissions:" + userID.String()
	data, err := r.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var permission []*entities.Permission
	if err := json.Unmarshal(data, &permission); err != nil {
		return nil, err
	}
	return permission, nil
}

func (r *Permissioncache) Delete(ctx context.Context, userID uuid.UUID) error {
	key := "user:permissions:" + userID.String()
	return r.client.Del(ctx, key).Err()
}
