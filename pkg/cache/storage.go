package cache

import "errors"

func MakeStorage(c Config) (RedisStorage, error) {
	c.Log.Info().Msg("configuring cache storage")
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
