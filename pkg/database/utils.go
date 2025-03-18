package database

import (
	"encoding/json"

	"github.com/apsamuel/brainiac/pkg/common"
)

/*
RetrieveConfig retrieves the latest active configuration record from the database
*/
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
	rows, err := storage.ExecuteQuery("SELECT * FROM config_data WHERE active = true ORDER BY created_at DESC LIMIT 1")
	if err != nil {
		return nil, err
	}

	var configData ConfigDataSchema
	// take the first record from the result set, decrypt the data field, and return it
	if len(rows) > 0 {
		record := rows[0]
		// get the data field from the record
		if recordData, ok := record["data"].(string); ok {
			data, err := common.DecryptWithAESGCM(recordData, []byte(aesKey))
			if err != nil {
				return nil, err
			}
			return data, nil
		}

		return json.Marshal(configData)
	}
	return nil, nil

}

/*
PushConfig pushes a new configuration record to the database
*/
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
	/* set all previous records (ping active = false) */
	if _, err := storage.ExecuteQuery("UPDATE config_data SET active = false WHERE active = true"); err != nil {
		return err
	}

	if _, err := storage.Save(record); err != nil {
		return err
	}

	config.Log.Info().Msgf("config record with ID %s saved successfully.", recordId)

	return nil
}

/*
GetObservers returns a map of observer channels for each marked observer in the configuration

- observers are marked by setting the `extension` key in the configuration to `true`
*/
func GetObservers(jsonConfig map[string]interface{}) map[string]chan Item {
	observerChannels := make(map[string]chan Item)
	for key := range jsonConfig {
		observerChannels[key] = make(chan Item)
	}
	return observerChannels
}
