package database

import (
	"encoding/json"
	"reflect"
	"time"
)

type TrainingDataSchema struct {
	Id             string    `gorm:"column:id;primaryKey" json:"id"`
	EmbeddingId    string    `gorm:"column:embedding_id" json:"embedding_id"`
	Embedding      []float64 `gorm:"column:embedding;type:bytes" json:"embedding"`
	EmbeddingModel string    `gorm:"column:embedding_model" json:"embedding_model"`
	Source         string    `gorm:"column:source" json:"source"`
	SourceURL      string    `gorm:"column:source_url" json:"source_url"`
	ChunksTotal    int       `gorm:"column:chunks_total" json:"chunks_total"`
	ChunksIndexed  int       `gorm:"column:chunks_indexed" json:"chunks_indexed"`
	Content        string    `gorm:"column:content" json:"content"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	IndexedAt      time.Time `gorm:"column:indexed_at" json:"indexed_at"`
	IsActive       bool      `gorm:"column:is_active" json:"is_active"`
	Category       string    `gorm:"column:category" json:"category"`
	Metadata       string    `gorm:"column:metadata" json:"metadata"`
}

func (t TrainingDataSchema) TableName() string {
	return "training_data"
}

func (t TrainingDataSchema) Schema() map[string]string {
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

func (t TrainingDataSchema) Columns() []string {
	var c []string
	v := reflect.ValueOf(t)
	for i := 0; i < v.NumField(); i++ {
		c = append(c, v.Type().Field(i).Name)
	}
	return c
}

func (t TrainingDataSchema) GetId() string {
	return t.Id
}

func (t TrainingDataSchema) String() string {
	jsonBytes, err := json.Marshal(t)
	if err != nil {
		return ""
	}

	return string(jsonBytes)
}
