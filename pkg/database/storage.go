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
		return newPostgresStorage[T](c, tableName)
	default:
		return nil
	}
}

func MakeStorage(c Config) (*Storage, error) {
	var storage Storage
	storage.Name = c.Options.Dataset
	storage.Type = c.Options.Engine
	storage.TrainingData = newStorage[TrainingDataSchema](c, "training_data")
	storage.ConfigData = newStorage[ConfigDataSchema](c, "config_data")
	return &storage, nil
}
