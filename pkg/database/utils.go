package database

import (
	"github.com/apsamuel/brainiac/pkg/common"
)

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
	config := Config{
		Options: Options{
			Dataset: configDatabase,
			Engine:  "postgres",
			Postgres: PostgresConfig{
				Host:        configHost,
				Port:        configPort,
				Username:    configUser,
				Password:    configPassword,
				DatasetName: configDatabase,
			},
		},
	}

	config.Log = common.GetLogger()

	configInterface := config.ToInterface()

	err := config.ConfigureFromInterface(configInterface)
	if err != nil {
		panic(err)
	}

	PostgresClient, err = NewPostgresClient(config.Options.Postgres)
	if err != nil {
		panic(err)
	}

	cipherText, err := common.EncryptWithAESGCM(data, []byte(aesKey))

	if err != nil {
		panic(err)
	}

	record := ConfigDataSchema{
		Id:        common.GetUUID(),
		Data:      cipherText,
		CreatedAt: common.GetTimeNowUTC(),
	}
	storage := NewStorage[ConfigDataSchema](config, "config_data")

	storage.Save(record)
	return nil
}
