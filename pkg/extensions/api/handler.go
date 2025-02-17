package api

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/apsamuel/brainiac/pkg/cache"
	"github.com/apsamuel/brainiac/pkg/common"
	"github.com/apsamuel/brainiac/pkg/database"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type MethodHandleFuncs struct {
	Methods    []string
	Public     string
	HandleFunc http.HandlerFunc
}
type Handler struct {
	Config                   *Config
	Router                   *mux.Router
	EndpointMethodHandlerMap map[string]*common.Route
	Storage                  *database.Storage
	Cache                    *cache.RedisStorage
	Observers                map[string]chan common.Item
	Templater                *template.Template
}

func (h *Handler) ConsumeEvents(eventChannel chan common.Item) error {
	for item := range eventChannel {
		h.Config.Log.Info().Msgf("Received event: %v", item)
	}
	return nil
}

func (h *Handler) ToEventChannel(item common.Item) {
	c, ok := h.Observers[item.Destination]
	if !ok {
		h.Config.Log.Error().Msg("Observer not found")
		return
	}

	c <- item
}

func (h *Handler) MakeRouter() error {
	h.Router = mux.NewRouter()
	functionMap := template.FuncMap{
		"foo": func() string {
			return "bar"
		},
	}
	h.Templater = template.Must(template.New("index").Funcs(functionMap).ParseGlob(fmt.Sprintf("%s/*.html", h.Config.Options.TemplateDir)))
	// h.EndpointMethodHandlerMap = make(map[string]MethodHandleFuncs)
	// _ = h.AddRoute("/health", []string{"GET"}, h.Health, "public")
	return nil
}

func (h *Handler) AddRoute(route *common.Route) error {

	if route == nil {
		return fmt.Errorf("route is nil")
	}

	h.Router.HandleFunc(route.Endpoint, route.Handler).Methods(route.Methods...)

	return nil
}

func (h *Handler) ListRoutes() []*common.Route {
	routes := []*common.Route{
		{
			Endpoint: "/health",
			Methods:  []string{"GET"},
			Handler:  h.Health,
			Auth:     "public",
		},
	}
	return routes
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
		Debug:                h.Config.Options.Debug,
		Logger:               h.Config.Log,
	}
	corsMiddleware := cors.New(corsOptions).Handler(h.Router)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", h.Config.Options.Host, h.Config.Options.Port), corsMiddleware); err != nil {
		h.Config.Log.Error().Err(err).Msg("failed to start server")
		return err
	}
	h.Config.Log.Info().Msg("server started")
	return nil
}
