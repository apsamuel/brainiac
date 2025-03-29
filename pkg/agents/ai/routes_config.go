package ai

import (
	"encoding/json"
	w3 "net/http"
)

func (agent *Agent) ConfigRequest(w w3.ResponseWriter, r *w3.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(agent.Config)
		_, err := w.Write(b)
		if err != nil {
			agent.Config.Log.Error().Err(err).Msg("failed to write response")
		}
		return
	} else {
		w.WriteHeader(w3.StatusMethodNotAllowed)
		_, err := w.Write([]byte(`{"error": "Method not allowed"}`))
		if err != nil {
			agent.Config.Log.Error().Err(err).Msg("failed to write error response")
		}
		return
	}

}
