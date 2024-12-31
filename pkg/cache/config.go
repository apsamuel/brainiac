package cache

import (
	"os"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Options Options `yaml:"cache"`
	Log     *zerolog.Logger
}

type Options struct {
	Engine string      `yaml:"engine"`
	Redis  RedisConfig `yaml:"redis"`
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
