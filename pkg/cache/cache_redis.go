package cache

import (
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password" json:"-"`
	Database int    `yaml:"database"`
}

type RedisStorage struct {
	client *redis.Client
}

func newRedisStorage(config RedisConfig) *RedisStorage {
	if config.Password == "" {
		return &RedisStorage{
			client: redis.NewClient(&redis.Options{
				Addr: config.Host + ":" + strconv.Itoa(config.Port),
				DB:   config.Database,
			}),
		}
	}

	return &RedisStorage{
		client: redis.NewClient(&redis.Options{
			Addr:     config.Host + ":" + strconv.Itoa(config.Port),
			Password: config.Password,
			DB:       config.Database,
		}),
	}
}
