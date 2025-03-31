package control

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
)

type ControlNode struct {
	Config *Config
	Log    *zerolog.Logger
	Router *mux.Router
}

func (node *ControlNode) Init() {
	// Initialize the control node
	node.Log = node.Config.Log
	node.Log.Info().Msg("control node initialized")
	node.Router = mux.NewRouter()
	node.startHttpServer()
}

func (node *ControlNode) startHttpServer() {
	// Start the server
	node.Log.Info().Msg("starting control node server")
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		Debug:            true,
		Logger:           node.Log,
		MaxAge:           0,
	}
	corsMiddleware := cors.New(corsOptions).Handler(node.Router)
	if err := http.ListenAndServe(
		node.Config.Options.Listen,
		corsMiddleware,
	); err != nil {
		node.Log.Error().Err(err).Msg("failed to start control node server")
	}
	node.Log.Info().Msg("control node server started")
	// Implement server start logic here
}

// func NewControlNode(config *Config) *ControlNode {
// 	// the control node controls agents, their configuration, persistence, and
// 	// the control node is responsible for managing the agents and their configuration
// 	controller := &ControlNode{
// 		Config: config,
// 		Log:    config.Log,
// 	}

// 	controller.Init()

// 	return controller
// }
