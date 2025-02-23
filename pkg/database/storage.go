package database

import (
	"github.com/apsamuel/brainiac/pkg/common"
)

type Storage struct {
	Name         string
	Type         string
	TrainingData common.Storer[TrainingDataSchema]
	ConfigData   common.Storer[ConfigDataSchema]
}

func newStorage[T any](c Config, tableName string) common.Storer[T] {
	switch c.Options.Engine {
	case "postgres":
		return NewPostgresStorage[T](c, tableName)
	default:
		return nil
	}
}

func MakeStorage(c Config) (*Storage, error) {
	var storage Storage
	storage.Name = c.Options.Dataset
	storage.Type = c.Options.Engine
	storage.TrainingData = newStorage[TrainingDataSchema](c, "training_data")
	// storage.ConfigData = newStorage[ConfigDataSchema](c, "config_data")
	return &storage, nil
}

func PushConfig(
	configHost string,
	configPort int,
	configDatabase string,
	configTable string,
	configUser string,
	configPassword string,
	data []byte,
	aesKey string,
	nonce string,
) error {
	postgresOptions := PostgresConfig{
		Host:        configHost,
		Port:        configPort,
		Username:    configUser,
		Password:    configPassword,
		DatasetName: configDatabase,
	}
	NewPostgresStorage[ConfigDataSchema](Config{
		Options: Options{
			Dataset:  configDatabase,
			Engine:   "postgres",
			Postgres: postgresOptions,
		},
	}, configTable)
	return nil
}
