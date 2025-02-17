package ai

import (
	"os"

	"gopkg.in/yaml.v2"
)

type ModelPersona struct {
	Model   string                 `yaml:"model"`
	System  string                 `yaml:"system"`
	Options map[string]interface{} `yaml:"options"`
}

type Options struct {
	Engine                string                  `yaml:"engine"`
	EmbeddingApiToken     string                  `yaml:"embedding_api_token"`
	EmbeddingURL          string                  `yaml:"embedding_url"`
	GenerateApiToken      string                  `yaml:"generate_api_token"`
	GenerateURL           string                  `yaml:"generate_url"`
	DefaultEmbeddingModel string                  `yaml:"embedding_model"`
	DefaultGenerateModel  string                  `yaml:"generate_model"`
	ModelPersonas         map[string]ModelPersona `yaml:"personas"`
}

type Config struct {
	Options Options `yaml:"ai"`
}

func (c *Config) Configure(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return err
	}
	return nil
}
