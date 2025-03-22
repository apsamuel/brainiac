package proxy

import (
	"github.com/apsamuel/brainiac/pkg/cache"
	"github.com/apsamuel/brainiac/pkg/database"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type Agent struct {
	Config    *Config
	Router    *mux.Router
	Storage   *database.Storage
	Log       zerolog.Logger
	Observers map[string]chan database.Item
	Cache     *cache.RedisStorage
}
