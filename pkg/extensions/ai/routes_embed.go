package ai

import (
	w3 "net/http"

	"github.com/apsamuel/brainiac/pkg/http"
)

type Response struct {
	Data string `json:"data"`
}

// EmbedRequest
func (h *Handler) EmbedRequest(w w3.ResponseWriter, r *w3.Request) {
	var request EmbedRequest
	_, err := http.UnmarshalRequestBody(&request, r)
	if err != nil {
		return
	}
}
