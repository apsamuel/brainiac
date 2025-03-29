package ai

import (
	"encoding/json"
	"html/template"
	"net/http"
)

const SelfPage = "aiHealth"

type HealthResponse struct {
	Status string `json:"status"`
}
type HealthPage struct {
	PageName   string         `json:"page_name" yaml:"page_name"`
	Config     *Config        `yaml:"api" json:"config"`
	JavaScript []template.JS  `json:"javascript" yaml:"javascript"`
	Style      []template.CSS `json:"style" yaml:"style"`
}

func (h *Agent) HealthRequest(w http.ResponseWriter, r *http.Request) {
	var data HealthResponse
	data.Status = "OK"
	b, _ := json.Marshal(data)
	if r.Header.Get("Content-Type") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(b)
		if err != nil {
			h.Config.Log.Error().Err(err).Msg("failed to write response")
		}
		return
	}
}
