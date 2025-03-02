package database

import (
	"fmt"

	"github.com/apsamuel/brainiac/pkg/common"
)

func RetrieveConfig(
	configHost string,
	configPort int,
	configDatabase string,
	configTable string,
	configUser string,
	configPassword string,
	aesKey string,
	nonce string,
) ([]byte, error) {
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
		return nil, err
	}

	PostgresClient, err = NewPostgresClient(config.Options.Postgres)
	if err != nil {
		return nil, err
	}

	storage := NewStorage[ConfigDataSchema](config, "config_data")

	// get the latest record by CreatedAt
	rows, err := storage.ExecuteQuery("SELECT * FROM config_data ORDER BY CreatedAt DESC LIMIT 1")
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		fmt.Println(row)
	}
	return nil, nil

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

	recordId := common.GetUUID()
	recordCreatedAt := common.GetTimeNowUTC()
	record := ConfigDataSchema{
		Id:        recordId,
		Data:      cipherText,
		CreatedAt: recordCreatedAt,
		Active:    true,
	}
	storage := NewStorage[ConfigDataSchema](config, "config_data")

	if _, err := storage.Save(record); err != nil {
		return err
	}

	config.Log.Info().Msgf("Record with ID %s saved successfully.", recordId)

	return nil
}
