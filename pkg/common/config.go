package common

import (
	"encoding/json"
	"os"

	// "github.com/apsamuel/brainiac/pkg/database"

	"gopkg.in/yaml.v2"
)

type ApiConfig struct {
	Origins     []string `json:"origins" yaml:"origins"`
	Host        string   `json:"host" yaml:"host"`
	Port        int      `json:"port" yaml:"port"`
	StaticDir   string   `json:"static_dir" yaml:"static_dir"`
	TemplateDir string   `json:"template_dir" yaml:"template_dir"`
	Debug       bool     `json:"debug" yaml:"debug"`
}

type PersonaConfig struct {
	Name         string                 `json:"name" yaml:"name"`
	SystemPrompt string                 `json:"system_prompt" yaml:"system_prompt"`
	Model        string                 `json:"model" yaml:"model"`
	ModelOptions map[string]interface{} `json:"model_options" yaml:"model_options"`
}

type ModelPersona struct {
	Model   string                 `yaml:"model" json:"model"`
	System  string                 `yaml:"system" json:"system"`
	Options map[string]interface{} `yaml:"ai" json:"ai"`
}

type AiConfig struct {
	Engine                string                  `yaml:"engine" json:"engine"`
	EmbeddingApiToken     string                  `yaml:"embedding_api_token" json:"embedding_api_token"`
	EmbeddingURL          string                  `yaml:"embedding_url" json:"embedding_url"`
	GenerateApiToken      string                  `yaml:"generate_api_token" json:"generate_api_token"`
	GenerateURL           string                  `yaml:"generate_url" json:"generate_url"`
	DefaultEmbeddingModel string                  `yaml:"embedding_model" json:"embedding_model"`
	DefaultGenerateModel  string                  `yaml:"generate_model" json:"generate_model"`
	ModelPersonas         map[string]ModelPersona `yaml:"personas" json:"personas"`
}
type RedisConfig struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Password string `json:"password" yaml:"password"`
}

type CacheConfig struct {
	Engine string      `json:"engine" yaml:"engine"`
	Redis  RedisConfig `json:"redis" yaml:"redis"`
}

type PostgresConfig struct {
	Host        string `json:"host" yaml:"host"`
	Port        int    `json:"port" yaml:"port"`
	Username    string `json:"username" yaml:"username"`
	Password    string `json:"password" yaml:"password"`
	DatasetName string `json:"dataset_name" yaml:"dataset_name"`
}

type SqliteConfig struct {
}

type DatabaseConfig struct {
	Engine   string         `json:"engine" yaml:"engine"`
	Postgres PostgresConfig `json:"postgres" yaml:"postgres"`
	Sqlite   SqliteConfig   `json:"sqlite" yaml:"sqlite"`
}

type ConfigInterface interface {
	GetOptions() map[string]interface{}
	ToInterface() map[string]interface{}
	FromInterface(map[string]interface{}) (Config, error)
	FromFile(string) (Config, error)
	String() string
	Publish() error
	Retrieve() error
}

type Config struct {
	Api      ApiConfig      `json:"api" yaml:"api"`
	Ai       AiConfig       `json:"ai" yaml:"ai"`
	Database DatabaseConfig `json:"database" yaml:"database"`
	Cache    CacheConfig    `json:"cache" yaml:"cache"`
}

/*
ToInterface converts a Config struct to an interface
*/
func (c *Config) ToInterface() (map[string]interface{}, error) {
	jsonConfig := make(map[string]interface{})
	jsonMap, err := json.Marshal(c)
	if err != nil {
		return jsonConfig, err
	}
	err = json.Unmarshal(jsonMap, &jsonConfig)
	if err != nil {
		return jsonConfig, err
	}
	return jsonConfig, nil
}

/*
FromInterface configures a Config struct from an interface with the internal structure of `map[string]any`

- jsonConfig: the interface with the internal structure of `map[string]any`
*/
func (c *Config) FromInterface(jsonConfig map[string]interface{}) (Config, error) {
	jsonMap, err := json.Marshal(jsonConfig)
	if err != nil {
		return Config{}, err
	}
	err = json.Unmarshal(jsonMap, &c)
	if err != nil {
		return Config{}, err
	}
	return *c, nil
}

/*
FromString configures a Config struct from a YAML encoded string

- yamlConfig: the YAML encoded string
*/
func (c *Config) FromString(yamlConfig string) (Config, error) {
	err := yaml.Unmarshal([]byte(yamlConfig), &c)
	if err != nil {
		return Config{}, err
	}
	return *c, nil
}

/*
FromBytes configures a Config struct from a YAML encoded byte slice

- yamlConfig: the YAML encoded byte slice
*/
func (c *Config) FromBytes(yamlConfig []byte) (Config, error) {
	err := yaml.Unmarshal(yamlConfig, &c)
	if err != nil {
		return Config{}, err
	}
	return *c, nil
}

/*
FromFile configures a Config struct from a YAML file

- filename: the path to the YAML configuration file
*/
func (c *Config) FromFile(filename string) (Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return Config{}, err
	}
	return *c, nil
}

/*
String converts a Config struct to a YAML encoded string
*/
func (c *Config) String() string {
	jsonConfig, err := c.ToInterface()
	if err != nil {
		return ""
	}
	jsonMap, err := json.Marshal(jsonConfig)
	if err != nil {
		return ""
	}
	return string(jsonMap)
}

/*
FromFileToStruct parses a YAML configuration and returns the provided type

- filename: the path to the YAML configuration file

- c: the type to return
*/
func FromFileToStruct(filename string, c any) (any, error) {
	if c == nil {
		c = make(map[string]interface{})
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func FromInterfaceToStruct(data map[string]interface{}, c any) (any, error) {
	if c == nil {
		c = make(map[string]interface{})
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(dataBytes, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
