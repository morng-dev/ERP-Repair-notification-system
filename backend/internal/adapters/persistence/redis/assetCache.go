package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/morng-dev/erp/internal/core/domain/entities"
	"github.com/morng-dev/erp/internal/core/domain/ports/cache"
	"github.com/redis/go-redis/v9"
)

type AseetCache struct {
	client *redis.Client
}

func NewAseetCache(client *redis.Client) cache.AssetCacheRepo {
	return &AseetCache{client: client}
}

func GetAssetKey(page, limit int) string {
	return fmt.Sprintf("assets:list:page:%d:limit:%d", page, limit)
}

func (r *AseetCache) SetAll(ctx context.Context, page, limit int, data *entities.AssetListCacheData) error {
	key := GetAssetKey(page, limit)
	value, err := json.Marshal(data)
	if err != nil {
		log.Printf("Redis SET Aseet ERROR: %v", err)
		return err
	}
	return r.client.Set(ctx, key, value, 5*time.Minute).Err()
}

func (r *AseetCache) GetAll(ctx context.Context, page, limit int) (*entities.AssetListCacheData, error) {
	key := GetAssetKey(page, limit)
	dataByte, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	var result entities.AssetListCacheData
	if err := json.Unmarshal(dataByte, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
