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
type Agent struct {
	Config                   *Config
	Router                   *mux.Router
	EndpointMethodHandlerMap map[string]*common.Route
	Storage                  *database.Storage
	Cache                    *cache.RedisStorage
	Observers                map[string]chan database.Item
	Templater                *template.Template
}

func (agent *Agent) ConsumeEvents(eventChannel chan database.Item) error {
	for item := range eventChannel {
		agent.Config.Log.Info().Msgf("received event: %v", item)
	}
	return nil
}

func (agent *Agent) ToEventChannel(item database.Item) {
	c, ok := agent.Observers[item.Destination]
	if !ok {
		agent.Config.Log.Error().Msg("observer not found")
		return
	}

	c <- item
}

func (agent *Agent) MakeRouter() error {
	agent.Router = mux.NewRouter()
	functionMap := template.FuncMap{
		"foo": func() string {
			return "bar"
		},
	}
	/*
	  attach allowed origins and configure CORS
	*/
	// corsOptions := handlers.AllowedOrigins(h.Config.Options.Origins)
	// headersOptions = handlers.AllowedHeaders([]string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"})
	// handler := handlers.CORS(corsOptions, headersOptions)(h.Router)
	agent.Config.Log.Info().Msgf("loading templates from %s", agent.Config.Options.TemplateDir)
	agent.Templater = template.Must(template.New("index").Funcs(functionMap).ParseGlob(fmt.Sprintf("%s/*.html", agent.Config.Options.TemplateDir)))
	return nil
}

func (agent *Agent) AddRoute(route *common.Route) error {

	if route == nil {
		return fmt.Errorf("route is nil")
	}

	agent.Router.HandleFunc(route.Endpoint, route.Handler).Methods(route.Methods...)

	return nil
}

func (agent *Agent) ListRoutes() []*common.Route {
	routes := []*common.Route{
		{
			Endpoint: "/health",
			Methods:  []string{"GET"},
			Handler:  agent.Health,
			Auth:     "public",
		},
	}
	return routes
}

func (agent *Agent) Serve() error {
	corsOptions := cors.Options{
		AllowedOrigins:  agent.Config.Options.Origins,
		AllowOriginFunc: nil,
		AllowedMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept-Encoding",
			"Accept-Language",
			"Accept",
			"Alt-Used",
			"Authorization",
			"Cache-Control",
			"Connection",
			"Content-Length",
			"Content-Type",
			"Cookie",
			"Host",
			"HX-Current-URL",
			"HX-Request",
			"Allow-Origin",
			"Priority",
			"Referer",
			"Sec-Fetch-Dest",
			"Sec-Fetch-Mode",
			"Sec-Fetch-Site",
			"Sec-Fetch-User",
			"Upgrade-Insecure-Requests",
			"User-Agent",
			"X-CSRF-Token",
			"X-Forwarded-For",
			"X-Forwarded-Proto",
			"X-Okta-User-Agent-Extended",
			"X-Okta-User-Agent",
			"X-Real-IP",
			"X-Requested-With",
			"*",
		},
		ExposedHeaders:       nil,
		MaxAge:               0,
		AllowCredentials:     true,
		AllowPrivateNetwork:  true,
		OptionsPassthrough:   false,
		OptionsSuccessStatus: 204,
		Debug:                agent.Config.Options.Debug,
		Logger:               agent.Config.Log,
	}
	corsMiddleware := cors.New(corsOptions).Handler(agent.Router)

	agent.Config.Log.Info().Msgf("starting %s server on %s:%d", SelfName, agent.Config.Options.Host, agent.Config.Options.Port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", agent.Config.Options.Host, agent.Config.Options.Port), corsMiddleware); err != nil {
		agent.Config.Log.Error().Err(err).Msg(fmt.Sprintf("failed to start %s server", SelfName))
		return err
	}
	agent.Config.Log.Info().Msg(fmt.Sprintf("%s server started", SelfName))
	return nil
}
