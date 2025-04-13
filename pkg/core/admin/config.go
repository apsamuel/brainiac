package control

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
)

const SelfName = "control"

type Options struct {
	Listen string `yaml:"listen" json:"listen"`
	Enable bool   `yaml:"enabled" json:"enabled"`
	Host   string `yaml:"host" json:"host"`
	Port   int    `yaml:"port" json:"port"`
}
type Config struct {
	Options Options         `yaml:"control" json:"control"`
	Log     *zerolog.Logger `yaml:"-" json:"-"`
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
