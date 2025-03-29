package ai

import (
	"encoding/json"
	w3 "net/http"

	"github.com/apsamuel/brainiac/pkg/http"
)

// type Response struct {
// 	Data string `json:"data"`
// }

func (agent *Agent) EmbedRequest(w w3.ResponseWriter, r *w3.Request) {
	if r.Method == "POST" {
		var request EmbedRequest
		_, err := http.UnmarshalRequestBody(&request, r)
		if err != nil {
			return
		}

		response, err := agent.Embed(request)
		if err != nil {
			agent.Config.Log.Error().Err(err).Msg("failed to embed")
			return
		}
		b, _ := json.Marshal(response)
		if r.Header.Get("Content-Type") == "application/json" {
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write(b)
			if err != nil {
				agent.Config.Log.Error().Err(err).Msg("failed to write response")
			}
			return
		}
	} else {
		w.WriteHeader(w3.StatusMethodNotAllowed)
		_, err := w.Write([]byte(`{"error": "Method not allowed"}`))
		if err != nil {
			agent.Config.Log.Error().Err(err).Msg("failed to write error response")
		}
		return
	}

}
