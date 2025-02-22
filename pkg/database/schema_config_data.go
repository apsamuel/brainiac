package database

import (
	"encoding/json"
	"reflect"
)

type ConfigDataSchema struct {
	Id      string `json:"id"`
	Data    string `json:"data"`
	Version string `json:"version"`
}

func (t ConfigDataSchema) TableName() string {
	return "config_data"
}

func (t ConfigDataSchema) Schema() map[string]string {
	var m = make(map[string]string)
	v := reflect.ValueOf(t)
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		if field.Type.String() == "[]float64" {
			m[field.Name] = "blob"
		}
		if field.Type.String() == "string" {
			m[field.Name] = "text"
		}
		if field.Type.String() == "int" {
			m[field.Name] = "integer"
		}
		if field.Type.String() == "bool" {
			m[field.Name] = "integer"
		}
		if field.Type.String() == "time.Time" {
			m[field.Name] = "date"
		}
	}
	return m
}

func (t ConfigDataSchema) Columns() []string {
	var c []string
	v := reflect.ValueOf(t)
	for i := 0; i < v.NumField(); i++ {
		c = append(c, v.Type().Field(i).Name)
	}
	return c
}

func (t ConfigDataSchema) GetId() string {
	return t.Id
}

func (t ConfigDataSchema) String() string {
	jsonBytes, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func (s *Storage) PushConfigData(data ConfigDataSchema) error {
	return s.ConfigData.Save(data)
}
