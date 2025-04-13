package control

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strconv"

	"github.com/apsamuel/brainiac/pkg/cache"
	"github.com/apsamuel/brainiac/pkg/common"
	"github.com/apsamuel/brainiac/pkg/database"
	"github.com/apsamuel/brainiac/pkg/logger"
	"github.com/gorilla/mux"
)

type Agent struct {
	Config    *Config
	Router    *mux.Router
	Log       logger.ZeroLogger
	Templater *template.Template
	Storage   *database.Storage
	Cache     *cache.RedisStorage
	Observers map[string]chan database.Item
	Channel   chan database.Item
}

func NewAgent(jsonConfig map[string]interface{}, logger logger.ZeroLogger) (*Agent, error) {
	config := &Config{}
	_, err := config.FromInterface(jsonConfig)
	if err != nil {
		return nil, err
	}
	// config.Log = logger

	agent := &Agent{
		Config:    config,
		Observers: make(map[string]chan database.Item),
		Channel:   make(chan database.Item),
		Log:       logger,
	}

	templateDir := filepath.Join(common.AppRoot, "pkg", "agents", SelfName, "templates")
	agent.Templater = template.Must(template.New("base").Funcs(template.FuncMap{}).ParseGlob(fmt.Sprintf("%s/*.html", templateDir)))
	agent.Router = mux.NewRouter()

	return agent, nil
}

func (agent *Agent) AddRoute(route *common.Route) error {
	if route == nil {
		return fmt.Errorf("route is nil")
	}

	//
	agent.Router.HandleFunc(route.Endpoint, route.Handler).Methods(route.Methods...)
	return nil
}

func (agent *Agent) ListRoutes() []*common.Route {

	routes := []*common.Route{}

	return routes
}

func (agent *Agent) Serve() error {
	agent.Config.Log.Info().Msgf("starting %s server on %s:%s", SelfName, agent.Config.Options.Host, strconv.Itoa(agent.Config.Options.Port))
	return nil
}
