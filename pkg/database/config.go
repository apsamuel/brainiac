package database

import (
	"os"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
)

type PostgresConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	DatasetName string `yaml:"dataset_name"`
}

type SqliteConfig struct {
	Filename    string `yaml:"filename"`
	DatasetName string `yaml:"dataset_name"`
}

type Options struct {
	Dataset  string         `yaml:"dataset"`
	Engine   string         `yaml:"engine"`
	Postgres PostgresConfig `yaml:"postgres"`
	Sqlite   SqliteConfig   `yaml:"sqlite"`
}

type Config struct {
	Options Options `yaml:"database"`
	Log     *zerolog.Logger
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

	if c.Options.Engine == "postgres" {
		c.Options.Postgres.DatasetName = c.Options.Dataset
		PostgresClient, err = NewPostgresClient((c.Options.Postgres))
		if err != nil {
			return err
		}
	}

	if c.Options.Engine == "sqlite" {
		c.Options.Sqlite.DatasetName = "main." + c.Options.Dataset
		SqliteClient, err = NewSqliteClient((c.Options.Sqlite))
		if err != nil {
			return err
		}
	}

	return nil
}
