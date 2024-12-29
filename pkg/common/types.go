package common

import (
	"context"
)

type Storer[T any] interface {
	// Get(key string) (T, error)
	// Set(key string, value T) error
	Retrieve(query string) ([]T, error)
	RetrieveById(id string) ([]T, error)
	VectorSearch(queryVector []float64) ([]T, error)
	ExecuteQuery(ctx context.Context, query string, args ...interface{}) ([]interface{}, error)
}

type Schema interface {
	String() string
	GetId() string
	TableName() string
	Columns() []string
}
