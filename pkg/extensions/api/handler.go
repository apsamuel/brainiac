package api

import (
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/apsamuel/brainiac/pkg/cache"
	"github.com/apsamuel/brainiac/pkg/common"
	"github.com/apsamuel/brainiac/pkg/database"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
)

type Handler struct {
	Config    *Config
	Router    *mux.Router
	Storage   *database.Storage
	Cache     *cache.RedisStorage
	Log       zerolog.Logger
	Observers map[string]chan common.Item
	Templater *template.Template
}

func (h *Handler) MakeRouter() error {
	h.Router = mux.NewRouter()
	h.Log = zerolog.New(os.Stdout).With().Timestamp().Logger()
	return nil
}

func (h *Handler) Serve() error {
	corsOptions := cors.Options{
		AllowedOrigins:       h.Config.Options.Origins,
		AllowOriginFunc:      nil,
		AllowedMethods:       []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:       []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:       nil,
		MaxAge:               0,
		AllowCredentials:     true,
		AllowPrivateNetwork:  true,
		OptionsPassthrough:   false,
		OptionsSuccessStatus: 204,
		Debug:                false,
		// Logger:               logger.Logger,
	}
	corsMiddleware := cors.New(corsOptions).Handler(h.Router)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", h.Config.Options.Host, h.Config.Options.Port), corsMiddleware); err != nil {
		h.Log.Error().Err(err).Msg("Failed to start server")
		return err
	}
	return nil
}
