package ai

import (
	"encoding/json"
	"fmt"

	"github.com/apsamuel/brainiac/pkg/cache"
	"github.com/apsamuel/brainiac/pkg/common"
	"github.com/apsamuel/brainiac/pkg/database"
	"github.com/apsamuel/brainiac/pkg/http"
)

type Agent struct {
	Config    *Config
	Channel   chan database.Item
	Observers map[string]chan database.Item
	Storage   *database.Storage
	Cache     *cache.RedisStorage
}

type GenerateResponse struct {
	Model              string    `json:"model"`
	CreatedAt          string    `json:"created_at"`
	Response           string    `json:"response"`
	Done               bool      `json:"done"`
	Context            []float64 `json:"context"`
	TotalDuration      float64   `json:"total_duration"`
	LoadDuration       float64   `json:"load_duration"`
	PromptEvalDuration float64   `json:"prompt_eval_duration"`
	EvalCount          float64   `json:"eval_count"`
	EvalDuration       float64   `json:"eval_duration"`
}

type GenerateResponseWrapper struct {
	Generate GenerateResponse         `json:"generate"`
	Stats    *http.TracedRequestStats `json:"stats"`
}

type EmbedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type EmbedResponse struct {
	Model      string    `json:"model"`
	Embeddings []float64 `json:"embeddings"`
}

type EmbedReponseWrapper struct {
	Embed EmbedResponse            `json:"embed"`
	Stats *http.TracedRequestStats `json:"stats"`
}

func makeHeaders(apiToken string) map[string]string {
	if apiToken == "" {
		return map[string]string{
			"Accept":       "application/json",
			"Content-Type": "application/json",
		}
	}
	return map[string]string{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", apiToken),
	}
}

func makeEmbedBody(input, modeName string) ([]byte, error) {
	var request EmbedRequest
	request.Input = input
	request.Model = modeName
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func makeGenerateBody(prompt, modelName, system string, context []float64) ([]byte, error) {
	var request GenerateRequest
	request.Prompt = prompt
	request.Model = modelName
	request.System = system
	request.Context = context
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type GenerateRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Stream  bool                   `json:"stream"`
	System  string                 `json:"system,omitempty"`
	Context []float64              `json:"context,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}

func (h *Agent) Consume(channel chan database.Item) error {
	for item := range channel {
		fmt.Printf("Received event: %v", item)
	}
	return nil
}

func (h *Agent) Generate(r GenerateRequest) (*GenerateResponseWrapper, error) {
	headers := makeHeaders("")
	var response GenerateResponse
	body, err := makeGenerateBody(r.Prompt, r.Model, r.System, r.Context)
	if err != nil {
		return nil, err
	}
	stats, err := Client.Query(
		h.Config.Options.GenerateURL,
		&response,
		headers,
		true,
		15000,
		body,
	)

	if err != nil {
		return nil, err
	}
	if stats.StatusCode != 200 {
		return nil, fmt.Errorf("status code %d", stats.StatusCode)
	}

	return &GenerateResponseWrapper{
		Generate: response,
		Stats:    stats,
	}, nil

}

func (h *Agent) Embed(
	r EmbedRequest,
) (*EmbedReponseWrapper, error) {
	headers := makeHeaders("")
	var response EmbedResponse
	body, err := makeEmbedBody(r.Input, r.Model)
	if err != nil {
		return nil, err
	}
	stats, err := Client.Query(
		h.Config.Options.EmbeddingURL,
		&response,
		headers,
		true,
		15000,
		body,
	)

	if err != nil {
		return nil, err
	}
	if stats.StatusCode != 200 {
		return nil, fmt.Errorf("status code %d", stats.StatusCode)
	}
	return &EmbedReponseWrapper{
		Embed: response,
		Stats: stats,
	}, nil

}

func (h *Agent) ListRoutes() []*common.Route {
	routes := []*common.Route{
		{
			Endpoint: "/ai/embed",
			Methods:  []string{"POST"},
			Handler:  h.EmbedRequest,
			Auth:     "public",
		},
	}
	return routes
}
