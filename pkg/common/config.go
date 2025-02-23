package common

import (
	"encoding/json"
	"os"

	// "github.com/apsamuel/brainiac/pkg/database"
	"gopkg.in/yaml.v2"
)

type PostgresConfigV2 struct {
	Host        string `json:"host" yaml:"host"`
	Port        int    `json:"port" yaml:"port"`
	Username    string `json:"username" yaml:"username"`
	Password    string `json:"password" yaml:"password"`
	DatasetName string `json:"dataset_name" yaml:"dataset_name"`
}

type SqliteConfigV2 struct {
}

type RedisConfigV2 struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Password string `json:"password" yaml:"password"`
}

type ApiConfigV2 struct {
	Origins     []string `json:"origins" yaml:"origins"`
	Host        string   `json:"host" yaml:"host"`
	Port        int      `json:"port" yaml:"port"`
	StaticDir   string   `json:"static_dir" yaml:"static_dir"`
	TemplateDir string   `json:"template_dir" yaml:"template_dir"`
	Debug       bool     `json:"debug" yaml:"debug"`
}

type PersonaConfigV2 struct {
	Name         string                 `json:"name" yaml:"name"`
	SystemPrompt string                 `json:"system_prompt" yaml:"system_prompt"`
	Model        string                 `json:"model" yaml:"model"`
	ModelOptions map[string]interface{} `json:"model_options" yaml:"model_options"`
}

type AiConfigV2 struct {
	Engine         string            `json:"engine" yaml:"engine"`
	EmbeddingUrl   string            `json:"embedding_url" yaml:"embedding_url"`
	EmbeddingModel string            `json:"embedding_model" yaml:"embedding_model"`
	GenerateUrl    string            `json:"generate_url" yaml:"generate_url"`
	GenerateModel  string            `json:"generate_model" yaml:"generate_model"`
	Personas       []PersonaConfigV2 `json:"personas" yaml:"personas"`
}

type CacheConfigV2 struct {
	Engine string        `json:"engine" yaml:"engine"`
	Redis  RedisConfigV2 `json:"redis" yaml:"redis"`
}

type DatabaseConfigV2 struct {
	Engine   string           `json:"engine" yaml:"engine"`
	Postgres PostgresConfigV2 `json:"postgres" yaml:"postgres"`
	Sqlite   SqliteConfigV2   `json:"sqlite" yaml:"sqlite"`
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
	Api      ApiConfigV2      `json:"api" yaml:"api"`
	Ai       AiConfigV2       `json:"ai" yaml:"ai"`
	Database DatabaseConfigV2 `json:"database" yaml:"database"`
	Cache    CacheConfigV2    `json:"cache" yaml:"cache"`
}

func (c *Config) GetOptions() map[string]interface{} {
	// return a list of defined option keys
	return map[string]interface{}{
		"api":      c.Api,
		"ai":       c.Ai,
		"database": c.Database,
		"cache":    c.Cache,
	}
}

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

func ParseConfig(filename string, c any) (any, error) {
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

// storage := database.NewPostgresStorage(database.PostgresConfig{
// 	Host:        configHost,
// 	Port:        configPort,
// 	Username:    configUser,
// 	Password:    configPassword,
// 	DatasetName: configDatabase,
// }, configTable)
// client, err := database.NewPostgresClient(database.PostgresConfig{
// 	Host:        configHost,
// 	Port:        configPort,
// 	Username:    configUser,
// 	Password:    configPassword,
// 	DatasetName: configDatabase,
// })
// if err != nil {
// 	return err
// }
// if err != nil {
// 	return err
// }
// return nil
// }
