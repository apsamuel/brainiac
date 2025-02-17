package common

import (
	"context"
	"net/http"
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

type Item struct {
	Source      string
	Destination string
	Attributes  []byte
	Value       Schema
}

type Route struct {
	Endpoint string
	Methods  []string
	Handler  http.HandlerFunc
	Auth     string
}
