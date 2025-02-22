package common

import (
	"os"

	"gopkg.in/yaml.v2"
)

func ParseConfig(filename string, c any) (any, error) {
	if c == nil {
		c = make(map[string]interface{})
	}

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
