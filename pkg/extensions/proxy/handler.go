package proxy

import (
	"github.com/apsamuel/brainiac/pkg/cache"
	"github.com/apsamuel/brainiac/pkg/common"
	"github.com/apsamuel/brainiac/pkg/database"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type Handler struct {
	Config    *Config
	Router    *mux.Router
	Storage   *database.Storage
	Log       zerolog.Logger
	Observers map[string]chan common.Item
	Cache     *cache.RedisStorage
}
