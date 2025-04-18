package database

import (
	"context"
	"database/sql"
)

// DBDriver defines the interface that all database drivers must implement
type DBDriver interface {
	// Connection management
	Connect(config DBConfig) error 
	Close() error
	Ping() error

	// Transaction management
	BeginTx(ctx context.Context) (*sql.Tx, error) 

	//Query execution
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row

	// Schema management
	InitializeSchema() error

	// Helper methods
	GetDialect() string
	TransformQuery(query string) string
}
