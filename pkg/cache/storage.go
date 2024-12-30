package cache

import "errors"

func MakeStorage(c Config) (RedisStorage, error) {
	switch c.Options.Engine {
	case "redis":
		storage, err := newRedisStorage(c.Options.Redis)
		if err != nil {
			return RedisStorage{}, err
		}
		return *storage, nil
	default:
		return RedisStorage{}, errors.New("unsupported cache engine")
	}
}
