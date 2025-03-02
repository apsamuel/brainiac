package database

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
)

type PostgresConfig struct {
	Host        string `yaml:"host" json:"host"`
	Port        int    `yaml:"port" json:"port"`
	Username    string `yaml:"username" json:"username"`
	Password    string `yaml:"password" json:"password"`
	DatasetName string `yaml:"dataset_name" json:"dataset_name"`
}

type SqliteConfig struct {
	Filename    string `yaml:"filename" json:"filename"`
	DatasetName string `yaml:"dataset_name" json:"dataset_name"`
}

type Options struct {
	Dataset  string         `yaml:"dataset" json:"dataset"`
	Engine   string         `yaml:"engine" json:"engine"`
	Postgres PostgresConfig `yaml:"postgres" json:"postgres"`
	Sqlite   SqliteConfig   `yaml:"sqlite" json:"sqlite"`
}

type Config struct {
	Options Options         `yaml:"database" json:"database"`
	Log     *zerolog.Logger `yaml:"-" json:"-"`
}

/*
ConfigureFromFile configures a Config struct from a YAML file
*/
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

/*
ConfigureFromInterface configures a Config struct from an interface `map[string]any`
*/
func (c *Config) ConfigureFromInterface(data map[string]interface{}) error {
	if databaseInterface, ok := data["database"].(map[string]any); ok {
		// I need to convert the databaseInterface of type `map[string]any` to a JSON encoded string
		// then unmarshal it to a Config struct

		dataBytes, err := json.Marshal(databaseInterface)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(dataBytes, c)
		if err != nil {
			return err
		}

		if c.Options.Engine == "postgres" {
			c.Options.Postgres.DatasetName = c.Options.Dataset
			PostgresClient, err = NewPostgresClient(c.Options.Postgres)
			if err != nil {
				return err
			}
		}

		if c.Options.Engine == "sqlite" {
			c.Options.Sqlite.DatasetName = "main." + c.Options.Dataset
			SqliteClient, err = NewSqliteClient(c.Options.Sqlite)
			if err != nil {
				return err
			}
		}
		return nil
	} else {
		return nil
	}
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
