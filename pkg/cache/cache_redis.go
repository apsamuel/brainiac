package cache

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func (r *RedisStorage) Set(ctx context.Context, key string, value any, expiration int) error {
	return r.client.Set(ctx, key, value, time.Duration(expiration)).Err()
}

func (r *RedisStorage) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return val, err
}

func (r *RedisStorage) GetAll(ctx context.Context, key string) (map[string]string, error) {
	val, err := r.client.HGetAll(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	return val, err
}

func (r *RedisStorage) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisStorage) Keys(ctx context.Context, pattern string) ([]string, error) {
	return r.client.Keys(ctx, pattern).Result()
}

func newRedisStorage(config RedisConfig) (*RedisStorage, error) {
	if config.Password == "" {
		storage := &RedisStorage{
			client: redis.NewClient(&redis.Options{
				Addr: config.Host + ":" + strconv.Itoa(config.Port),
				DB:   config.Database,
			}),
		}
		return storage, nil
	}

	storage := &RedisStorage{
		client: redis.NewClient(&redis.Options{
			Addr:     config.Host + ":" + strconv.Itoa(config.Port),
			Password: config.Password,
			DB:       config.Database,
		}),
	}
	return storage, nil
}
