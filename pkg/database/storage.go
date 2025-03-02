package database

type Storage struct {
	Name         string
	Type         string
	TrainingData Storer[TrainingDataSchema]
	ConfigData   Storer[ConfigDataSchema]
}

func NewStorage[T any](c Config, tableName string) Storer[T] {
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
	storage.TrainingData = NewStorage[TrainingDataSchema](c, "training_data")
	/*
		  This line was commented out because ConfigDataSchema is created once during the bootstrapping process and may be updated later
			instead of being created here, we create it when brainiac is provided the config argument on the command line
		  storage.ConfigData = newStorage[ConfigDataSchema](c, "config_data")
	*/

	return &storage, nil
}

// func
