package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

const SelfPage = "health"

type Response struct {
	Data string `json:"data"`
}

type HealthPage struct {
	PageName string  `json:"page_name" yaml:"page_name"`
	Config   *Config `yaml:"api" json:"config"`
	// Handler    *Handler       `json:"handler" yaml:"handler"`
	JavaScript []template.JS  `json:"javascript" yaml:"javascript"`
	Style      []template.CSS `json:"style" yaml:"style"`
	// Template   string         `json:"template" yaml:"template"`
}

func (h *Agent) Health(w http.ResponseWriter, r *http.Request) {
	var data Response
	data.Data = "OK"
	b, _ := json.Marshal(data)
	if r.Header.Get("Content-Type") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(b)
		if err != nil {
			h.Config.Log.Error().Err(err).Msg("failed to write response")
		}
		return
	} else {
		w.Header().Set("Content-Type", "text/html")
		err := h.Templater.ExecuteTemplate(w, fmt.Sprintf("%s.html", SelfPage), HealthPage{
			PageName:   "Health",
			Config:     h.Config,
			JavaScript: nil,
			Style:      nil,
			// Template:   "health",
			// Handler:  h,
		})
		if err != nil {
			h.Config.Log.Error().Err(err).Msg("failed to write response")
		}
	}
	// w.Header().Set("Content-Type", "application/json")
	// _, err := w.Write(b)
	// if err != nil {
	// 	h.Config.Log.Error().Err(err).Msg("failed to write response")
	// }
}
