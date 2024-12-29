package cache

func MakeStorage(c Config) RedisStorage {
	switch c.Options.Engine {
	case "redis":
		storage := newRedisStorage(c.Options.Redis)
		return *storage
	default:
		return RedisStorage{}
	}
}
