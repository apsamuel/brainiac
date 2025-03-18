package cache

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Options Options         `yaml:"cache" json:"cache"`
	Log     *zerolog.Logger `yaml:"-" json:"-"`
}

type Options struct {
	Engine string      `yaml:"engine" json:"engine"`
	Redis  RedisConfig `yaml:"redis" json:"redis"`
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
