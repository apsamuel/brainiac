package ai

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
)

const SelfName = "ai"

type ModelPersona struct {
	Model   string                 `yaml:"model" json:"model"`
	System  string                 `yaml:"system" json:"system"`
	Options map[string]interface{} `yaml:"ai" json:"ai"`
}

type Options struct {
	Engine                string                  `yaml:"engine" json:"engine"`
	EmbeddingApiToken     string                  `yaml:"embedding_api_token" json:"embedding_api_token"`
	EmbeddingURL          string                  `yaml:"embedding_url" json:"embedding_url"`
	GenerateApiToken      string                  `yaml:"generate_api_token" json:"generate_api_token"`
	GenerateURL           string                  `yaml:"generate_url" json:"generate_url"`
	DefaultEmbeddingModel string                  `yaml:"embedding_model" json:"embedding_model"`
	DefaultGenerateModel  string                  `yaml:"generate_model" json:"generate_model"`
	ModelPersonas         map[string]ModelPersona `yaml:"personas" json:"personas"`
}

type Config struct {
	Options Options         `yaml:"ai" json:"ai"`
	Log     *zerolog.Logger `yaml:"-" json:"-"`
}

func (c *Config) ConfigureFromFile(filename string) error {
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

/*
ConfigureFromInterface configures a Config struct from an interface `map[string]any`
*/
func (c *Config) ConfigureFromInterface(data map[string]interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(dataBytes, c)
	if err != nil {
		return err
	}
	return nil
}

/*
String converts a Config struct to a JSON encoded string
*/
func (c *Config) String() string {
	// returns JSON encoded string
	data, err := yaml.Marshal(c)
	if err != nil {
		return ""
	}
	return string(data)
}

/*
ToInterface converts a Config struct to an interface
*/
func (c *Config) ToInterface() map[string]interface{} {
	interfaceData := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(c.String()), &interfaceData)
	if err != nil {
		return nil
	}
	return interfaceData
}
