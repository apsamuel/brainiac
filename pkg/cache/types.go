package cache

import "github.com/redis/go-redis/v9"

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password" json:"-"`
	Database int    `yaml:"database"`
}

type RedisStorage struct {
	client *redis.Client
}
