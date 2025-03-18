package database

import "gorm.io/gorm"

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

/*
Storer is an interface that defines the methods that a storage system should implement.
*/
type Storer[T any] interface {
	Save(data T) (*gorm.DB, error)
	Retrieve(query string) ([]T, error)
	RetrieveById(id string) ([]T, error)
	VectorSearch(queryVector []float64) ([]T, error)
	ExecuteQuery(query string, args ...interface{}) ([]map[string]interface{}, error)
}
