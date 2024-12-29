package common

import (
	"os"

	"gopkg.in/yaml.v2"
)

func ReadConfigFile(filename string, c any) (any, error) {

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
