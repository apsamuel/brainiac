package api

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Options struct {
	Origins     []string `yaml:"origins"`
	Host        string   `yaml:"host"`
	Port        int      `yaml:"port"`
	StaticDir   string   `yaml:"static_dir"`
	TemplateDir string   `yaml:"template_dir"`
}

type Config struct {
	Options Options `yaml:"api"`
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
