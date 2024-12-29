package database

import (
	"encoding/json"
	"reflect"
	"time"
)

type TrainingDataSchema struct {
	ID             string    `gorm:"column:ID" json:"id"`
	EmbeddingId    string    `gorm:"column:EmbeddingId" json:"embedding_id"`
	Embedding      []float64 `gorm:"column:Embedding;type:bytes" json:"embedding"`
	EmbeddingModel string    `gorm:"column:EmbeddingModel" json:"embedding_model"`
	Source         string    `gorm:"column:Source" json:"source"`
	SourceURL      string    `gorm:"column:SourceURL" json:"source_url"`
	ChunksTotal    int       `gorm:"column:ChunksTotal" json:"chunks_total"`
	ChunksIndexed  int       `gorm:"column:ChunksIndexed" json:"chunks_indexed"`
	Content        string    `gorm:"column:Content" json:"content"`
	CreatedAt      time.Time `gorm:"column:CreatedAt" json:"created_at"`
	IndexedAt      time.Time `gorm:"column:IndexedAt" json:"indexed_at"`
	IsActive       bool      `gorm:"column:IsActive" json:"is_active"`
	Category       string    `gorm:"column:Category" json:"category"`
	Metadata       string    `gorm:"column:Metadata" json:"metadata"`
}

func (t TrainingDataSchema) TableName() string {
	return "training_data"
}

func (t TrainingDataSchema) Columns() []string {
	var c []string
	// return []string{
	// 	"ID",
	// 	"EmbeddingId",
	// 	"Embedding",
	// 	"EmbeddingModel",
	// 	"Source",
	// 	"SourceURL",
	// 	"ChunksTotal",
	// 	"ChunksIndexed",
	// 	"Content",
	// 	"CreatedAt",
	// 	"IndexedAt",
	// 	"IsActive",
	// 	"Category",
	// 	"Metadata",
	// }
	// use reflection to get the columns
	v := reflect.ValueOf(t)
	for i := 0; i < v.NumField(); i++ {
		c = append(c, v.Type().Field(i).Name)
	}
	return c
}

func (t TrainingDataSchema) GetId() string {
	return t.ID
}

func (t TrainingDataSchema) String() string {
	jsonBytes, err := json.Marshal(t)
	if err != nil {
		return ""
	}

	return string(jsonBytes)
}
