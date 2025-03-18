package database

import (
	"encoding/json"
	"reflect"
	"time"
)

type ConfigDataSchema struct {
	// Id should be the primary key
	Id        string    `json:"id" gorm:"column:id;primaryKey"`
	Data      string    `json:"data" gorm:"column:data"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	Active    bool      `json:"active" gorm:"column:active"`
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
	return ""
}

func (t ConfigDataSchema) String() string {
	jsonBytes, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}
