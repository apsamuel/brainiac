package cache

import "github.com/redis/go-redis/v9"

type RedisConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Password string `yaml:"password" json:"-"`
	Database int    `yaml:"database"`
}

type RedisStorage struct {
	client *redis.Client
}
