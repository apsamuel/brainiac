package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var PostgresClient *gorm.DB

type PostgresStore[T any] struct {
	datasetName string
	tableName   string
}

func (s *PostgresStore[T]) Retrieve(query string) ([]T, error) {
	return nil, nil
}

func (s *PostgresStore[T]) RetrieveById(id string) ([]T, error) {
	return nil, nil
}

func (s *PostgresStore[T]) VectorSearch(queryVector []float64) ([]T, error) {
	return nil, nil
}

func (s *PostgresStore[T]) ExecuteQuery(query string, args ...interface{}) ([]interface{}, error) {
	rows, err := PostgresClient.Raw(query, args).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []interface{}
	for rows.Next() {
		var result interface{}
		if err := rows.Scan(&result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (s *PostgresStore[T]) Save(data T) (*gorm.DB, error) {
	result := PostgresClient.Save(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return result, nil
}

type Float64Slice []float64

// Value marshals Float64Slice to JSON for storage in the database
func (f Float64Slice) Value() (driver.Value, error) {
	return json.Marshal(f)
}

// Scan unmarshals JSON data to Float64Slice
func (f *Float64Slice) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &f)
}

// wrapFloat64SliceFields wraps []float64 fields with Float64Slice in the given struct
func wrapFloat64SliceFields(data interface{}) interface{} {
	v := reflect.ValueOf(data)
	t := v.Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	wrappedStruct := reflect.New(t).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type == reflect.TypeOf([]float64{}) {
			wrappedStruct.Field(i).Set(reflect.ValueOf(Float64Slice(v.Field(i).Interface().([]float64))))
		} else {
			wrappedStruct.Field(i).Set(v.Field(i))
		}
	}
	return wrappedStruct.Addr().Interface()
}

func checkPostgresTableExists(data any) bool {
	migrator := PostgresClient.Migrator()
	return migrator.HasTable(data)
}

func getSchema(data interface{}) map[string]string {
	m := make(map[string]string)
	v := reflect.ValueOf(data)
	t := v.Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.String() == "database.Float64Slice" {
			m[field.Name] = "blob"
		}
		if field.Type.String() == "db.PromptSchemaModelOptions" {
			m[field.Name] = "text"
		}
		if field.Type.String() == "time.Time" {
			m[field.Name] = "date"
		}
		if field.Type.String() == "float64" {
			m[field.Name] = "real"
		}
		if field.Type.String() == "string" {
			m[field.Name] = "text"
		}
		if field.Type.String() == "int" {
			m[field.Name] = "integer"
		}
		if field.Type.String() == "bool" {
			m[field.Name] = "boolean"
		}
	}
	return m
}

func createPostgresTable(c Config, data interface{}) error {
	tableName := data.(Schema).TableName()
	wrappedData := wrapFloat64SliceFields(data)

	if !checkPostgresTableExists(wrappedData) {
		c.Log.Info().Msg("table does not exist, creating table")
		schema := getSchema(data)
		var columns []string
		for k, v := range schema {
			c.Log.Info().Msg(fmt.Sprintf("%s %s", k, v))
			columns = append(columns, fmt.Sprintf("\"%s\" %s", k, v))
		}

		statement := fmt.Sprintf("CREATE TABLE %s (%s)", tableName, strings.Join(columns, ","))
		c.Log.Info().Msg(statement)
		result := PostgresClient.Exec(statement)
		if result.Error != nil {
			return result.Error
		}
		c.Log.Info().Msg("table created")
	} else {
		c.Log.Info().Msg("table already exists")
	}
	return nil
}

func buildPostgresDSN(config PostgresConfig) string {
	dsn := "host=" + config.Host +
		" port=" + strconv.Itoa(config.Port) +
		" user=" + config.Username +
		" password=" + config.Password +
		" dbname=" + config.DatasetName +
		" sslmode=disable"
	return dsn
}

func NewPostgresClient(config PostgresConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(buildPostgresDSN(config)), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewPostgresStorage[T any](config Config, tableName string) Storer[T] {
	var schema T

	if err := createPostgresTable(config, schema); err != nil {
		config.Log.Error().Msgf("error creating table %s", tableName)
	}
	s := new(PostgresStore[T])
	s.datasetName = config.Options.Dataset
	s.tableName = tableName
	return s
}
