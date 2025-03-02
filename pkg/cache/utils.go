package cache

import (
	"context"
)

func PushConfig(
	configHost string,
	configPort int,
	configKey string,
	configData []byte,
) error {
	redisOptions := RedisConfig{
		Host:     configHost,
		Port:     configPort,
		Database: 0,
	}
	postgresOptions := Options{
		Engine: "redis",
		Redis:  redisOptions,
	}
	config := Config{
		Options: postgresOptions,
	}
	storage, err := newRedisStorage(config.Options.Redis)
	if err != nil {
		return err
	}
	err = storage.Set(context.Background(), configKey, configData, 0)
	if err != nil {
		return err
	}
	return nil
}
