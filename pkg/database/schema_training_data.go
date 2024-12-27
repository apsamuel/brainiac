package database

import (
	"time"
)

type TrainingDataSchema struct {
	ID             int       `gorm:"column:ID"`
	EmbeddingId    string    `gorm:"column:EmbeddingId"`
	Embedding      []float64 `gorm:"column:Embedding:type:bytes"`
	EmbeddingModel string    `gorm:"column:EmbeddingModel"`
	Source         string    `gorm:"column:Source"`
	SourceURL      string    `gorm:"column:SourceURL"`
	ChunksTotal    int       `gorm:"column:ChunksTotal"`
	ChunksIndexed  int       `gorm:"column:ChunksIndexed"`
	Content        string    `gorm:"column:Content"`
	CreatedAt      time.Time `gorm:"column:CreatedAt"`
	IndexedAt      time.Time `gorm:"column:IndexedAt"`
	IsActive       bool      `gorm:"column:IsActive"`
	Category       string    `gorm:"column:Category"`
	Metadata       string    `gorm:"column:Metadata"`
}
