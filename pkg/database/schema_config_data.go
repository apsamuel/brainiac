package database

import (
	"encoding/json"
	"reflect"
	"time"
)

type ConfigDataSchema struct {
	Id        string    `json:"id" gorm:"column:Id"`
	Data      string    `json:"data" gorm:"column:Data"`
	CreatedAt time.Time `json:"created_at" gorm:"column:CreatedAt"`
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
