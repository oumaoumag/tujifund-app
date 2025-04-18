package database

import (
	"context"
	"database/sql"
)

// BaseDriver provides common implementations for the DBDriver interface
type BaseDriver struct {
	db *sql.DB
}

// Close closes the database connection
func (d *BaseDriver) Close() error {
	return d.db.Close()
}

// Ping checks the database connection
func (d *BaseDriver) Ping() error {
	return d.db.Ping()
}

// BeginTx starts the a transaction
func (d *BaseDriver) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return d.db.BeginTx(ctx, nil)
}

// Exec executes a query without returning any rows
func (d *BaseDriver) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.db.Exec(query, args...)
}

// QueryRow executes a query that return a single row
func (d *BaseDriver) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.db.QueryRow(query, args...)
} 