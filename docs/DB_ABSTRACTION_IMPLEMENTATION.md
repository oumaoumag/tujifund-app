# Database Abstraction Layer Implementation Guide

This document provides a detailed guide for implementing the database abstraction layer that will allow the TujiFund/ChamaVault application to work with both SQLite and PostgreSQL databases.

## Overview

The database abstraction layer will:
1. Provide a unified interface for database operations
2. Handle database-specific SQL dialects
3. Manage connections appropriately for each database type
4. Allow runtime switching between database systems

## Implementation Steps

### 1. Database Interface Definition

Create a common interface that defines all database operations:

```go
// backend/database/interface.go
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
    
    // Query execution
    Exec(query string, args ...interface{}) (sql.Result, error)
    Query(query string, args ...interface{}) (*sql.Rows, error)
    QueryRow(query string, args ...interface{}) *sql.Row
    
    // Schema management
    InitializeSchema() error
    
    // Helper methods
    GetDialect() string
    TransformQuery(query string) string
}
```

### 2. Enhanced Configuration Structure

Update the configuration structure to support both database systems:

```go
// backend/database/config.go
package database

// DBConfig holds database configuration
type DBConfig struct {
    // Common settings
    Driver string // "sqlite" or "postgres"
    
    // SQLite specific
    SQLitePath string
    
    // PostgreSQL specific
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
    SSLMode  string
    
    // Connection pool settings
    MaxOpenConns int
    MaxIdleConns int
    ConnMaxLifetime int // in seconds
}

// NewSQLiteConfig creates a SQLite configuration
func NewSQLiteConfig(path string) DBConfig {
    return DBConfig{
        Driver:     "sqlite",
        SQLitePath: path,
        // Default connection pool settings
        MaxOpenConns: 1, // SQLite supports only one writer
        MaxIdleConns: 2,
        ConnMaxLifetime: 3600,
    }
}

// NewPostgresConfig creates a PostgreSQL configuration
func NewPostgresConfig(host string, port int, user, password, dbname, sslmode string) DBConfig {
    return DBConfig{
        Driver:   "postgres",
        Host:     host,
        Port:     port,
        User:     user,
        Password: password,
        DBName:   dbname,
        SSLMode:  sslmode,
        // Default connection pool settings
        MaxOpenConns: 10,
        MaxIdleConns: 5,
        ConnMaxLifetime: 3600,
    }
}
```

### 3. SQLite Driver Implementation

Implement the SQLite-specific driver:

```go
// backend/database/sqlite_driver.go
package database

import (
    "context"
    "database/sql"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    
    _ "modernc.org/sqlite" // Modern SQLite driver
)

// SQLiteDriver implements the DBDriver interface for SQLite
type SQLiteDriver struct {
    db   *sql.DB
    conf DBConfig
}

// Connect establishes a connection to the SQLite database
func (d *SQLiteDriver) Connect(conf DBConfig) error {
    // Create the database directory if it doesn't exist
    if err := os.MkdirAll(filepath.Dir(conf.SQLitePath), 0o755); err != nil {
        return fmt.Errorf("failed to create database directory: %w", err)
    }

    dsn := fmt.Sprintf("file:%s?cache=shared&_journal_mode=WAL", conf.SQLitePath)
    db, err := sql.Open("sqlite", dsn)
    if err != nil {
        return fmt.Errorf("failed to connect to SQLite database: %w", err)
    }

    // Test the connection
    if err = db.Ping(); err != nil {
        return fmt.Errorf("failed to ping SQLite database: %w", err)
    }

    // Enable foreign key support
    if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
        return fmt.Errorf("failed to enable foreign keys: %w", err)
    }

    // Set connection pool settings
    db.SetMaxOpenConns(conf.MaxOpenConns)
    db.SetMaxIdleConns(conf.MaxIdleConns)

    d.db = db
    d.conf = conf
    return nil
}

// Additional methods implementing the DBDriver interface...

// GetDialect returns the SQL dialect name
func (d *SQLiteDriver) GetDialect() string {
    return "sqlite"
}

// TransformQuery converts a generic SQL query to SQLite syntax
func (d *SQLiteDriver) TransformQuery(query string) string {
    // Handle any SQLite-specific transformations
    return query
}
```

### 4. PostgreSQL Driver Implementation

Implement the PostgreSQL-specific driver:

```go
// backend/database/postgres_driver.go
package database

import (
    "context"
    "database/sql"
    "fmt"
    "time"
    
    _ "github.com/lib/pq" // PostgreSQL driver
)

// PostgresDriver implements the DBDriver interface for PostgreSQL
type PostgresDriver struct {
    db   *sql.DB
    conf DBConfig
}

// Connect establishes a connection to the PostgreSQL database
func (d *PostgresDriver) Connect(conf DBConfig) error {
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        conf.Host, conf.Port, conf.User, conf.Password, conf.DBName, conf.SSLMode)
    
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return fmt.Errorf("failed to connect to PostgreSQL database: %w", err)
    }

    // Test the connection
    if err = db.Ping(); err != nil {
        return fmt.Errorf("failed to ping PostgreSQL database: %w", err)
    }

    // Set connection pool settings
    db.SetMaxOpenConns(conf.MaxOpenConns)
    db.SetMaxIdleConns(conf.MaxIdleConns)
    db.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetime) * time.Second)

    d.db = db
    d.conf = conf
    return nil
}

// Additional methods implementing the DBDriver interface...

// GetDialect returns the SQL dialect name
func (d *PostgresDriver) GetDialect() string {
    return "postgres"
}

// TransformQuery converts a generic SQL query to PostgreSQL syntax
func (d *PostgresDriver) TransformQuery(query string) string {
    // Convert SQLite placeholders (?) to PostgreSQL placeholders ($1, $2, etc.)
    parts := strings.Split(query, "?")
    if len(parts) <= 1 {
        return query
    }
    
    result := parts[0]
    for i := 1; i < len(parts); i++ {
        result += fmt.Sprintf("$%d", i) + parts[i]
    }
    
    // Handle other PostgreSQL-specific transformations
    return result
}
```

### 5. Factory Function

Create a factory function to return the appropriate driver:

```go
// backend/database/factory.go
package database

import (
    "fmt"
)

// NewDriver creates a new database driver based on the configuration
func NewDriver(conf DBConfig) (DBDriver, error) {
    switch conf.Driver {
    case "sqlite":
        driver := &SQLiteDriver{}
        if err := driver.Connect(conf); err != nil {
            return nil, err
        }
        return driver, nil
        
    case "postgres":
        driver := &PostgresDriver{}
        if err := driver.Connect(conf); err != nil {
            return nil, err
        }
        return driver, nil
        
    default:
        return nil, fmt.Errorf("unsupported database driver: %s", conf.Driver)
    }
}
```

### 6. Schema Management

Create database-specific schema creation logic:

```go
// backend/database/schema.go
package database

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

// SchemaManager handles database schema creation and updates
type SchemaManager struct {
    driver DBDriver
}

// NewSchemaManager creates a new schema manager
func NewSchemaManager(driver DBDriver) *SchemaManager {
    return &SchemaManager{
        driver: driver,
    }
}

// InitializeSchema creates the database schema
func (sm *SchemaManager) InitializeSchema() error {
    // Load the appropriate schema file based on the dialect
    dialect := sm.driver.GetDialect()
    schemaPath := filepath.Join("database", fmt.Sprintf("schema_%s.sql", dialect))
    
    // If dialect-specific schema doesn't exist, fall back to the generic one
    if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
        schemaPath = filepath.Join("database", "database_schema.sql")
    }
    
    // Read the schema file
    schema, err := os.ReadFile(schemaPath)
    if err != nil {
        return fmt.Errorf("failed to read schema file: %w", err)
    }
    
    // Split the schema into individual statements
    statements := strings.Split(string(schema), ";")
    
    // Execute each statement
    for _, stmt := range statements {
        stmt = strings.TrimSpace(stmt)
        if stmt == "" {
            continue
        }
        
        // Transform the query for the specific dialect
        stmt = sm.driver.TransformQuery(stmt)
        
        if _, err := sm.driver.Exec(stmt); err != nil {
            // Ignore "already exists" errors
            if !strings.Contains(err.Error(), "already exists") {
                return fmt.Errorf("failed to execute schema statement: %w", err)
            }
        }
    }
    
    return nil
}
```

### 7. Update Main Application

Update the main application to use the abstraction layer:

```go
// backend/main.go (partial update)
import (
    "log"
    "os"
    "strconv"
    
    "tujifund-app/backend/database"
)

func main() {
    // Determine which database to use from environment or config
    dbDriver := os.Getenv("DB_DRIVER")
    if dbDriver == "" {
        dbDriver = "sqlite" // Default to SQLite
    }
    
    var config database.DBConfig
    
    if dbDriver == "sqlite" {
        // Configure SQLite
        dbPath := os.Getenv("DB_PATH")
        if dbPath == "" {
            dbPath = "data/tujifund.db" // Default path
        }
        config = database.NewSQLiteConfig(dbPath)
    } else if dbDriver == "postgres" {
        // Configure PostgreSQL
        host := os.Getenv("DB_HOST")
        if host == "" {
            host = "localhost"
        }
        
        portStr := os.Getenv("DB_PORT")
        port, err := strconv.Atoi(portStr)
        if err != nil {
            port = 5432 // Default PostgreSQL port
        }
        
        user := os.Getenv("DB_USER")
        password := os.Getenv("DB_PASSWORD")
        dbName := os.Getenv("DB_NAME")
        sslMode := os.Getenv("DB_SSLMODE")
        if sslMode == "" {
            sslMode = "disable" // Default to disable for development
        }
        
        config = database.NewPostgresConfig(host, port, user, password, dbName, sslMode)
    } else {
        log.Fatalf("Unsupported database driver: %s", dbDriver)
    }
    
    // Create database driver
    driver, err := database.NewDriver(config)
    if err != nil {
        log.Fatalf("Failed to create database driver: %v", err)
    }
    
    // Initialize schema
    schemaManager := database.NewSchemaManager(driver)
    if err := schemaManager.InitializeSchema(); err != nil {
        log.Printf("Failed to initialize database schema: %v", err)
        log.Println("Attempting to continue with existing schema...")
    }
    
    // Continue with application setup...
}
```

## Testing the Abstraction Layer

Create tests to verify the abstraction layer works with both database systems:

```go
// backend/database/database_test.go
package database

import (
    "os"
    "testing"
)

func TestSQLiteDriver(t *testing.T) {
    // Create a temporary SQLite database
    tempFile := "test_sqlite.db"
    defer os.Remove(tempFile)
    
    config := NewSQLiteConfig(tempFile)
    driver, err := NewDriver(config)
    if err != nil {
        t.Fatalf("Failed to create SQLite driver: %v", err)
    }
    
    // Test basic operations
    // ...
}

func TestPostgresDriver(t *testing.T) {
    // Skip if PostgreSQL environment variables are not set
    if os.Getenv("TEST_PG_HOST") == "" {
        t.Skip("Skipping PostgreSQL tests: environment variables not set")
    }
    
    config := NewPostgresConfig(
        os.Getenv("TEST_PG_HOST"),
        5432, // Default port
        os.Getenv("TEST_PG_USER"),
        os.Getenv("TEST_PG_PASSWORD"),
        os.Getenv("TEST_PG_DBNAME"),
        "disable",
    )
    
    driver, err := NewDriver(config)
    if err != nil {
        t.Fatalf("Failed to create PostgreSQL driver: %v", err)
    }
    
    // Test basic operations
    // ...
}
```

## Migration Scripts

Create scripts to migrate data from SQLite to PostgreSQL:

```go
// backend/database/migrations/scripts/migrate.go
package main

import (
    "database/sql"
    "flag"
    "fmt"
    "log"
    "os"
    
    "tujifund-app/backend/database"
    
    _ "github.com/lib/pq"
    _ "modernc.org/sqlite"
)

func main() {
    // Parse command-line flags
    sqlitePath := flag.String("sqlite", "", "Path to SQLite database")
    pgHost := flag.String("pghost", "localhost", "PostgreSQL host")
    pgPort := flag.Int("pgport", 5432, "PostgreSQL port")
    pgUser := flag.String("pguser", "", "PostgreSQL user")
    pgPassword := flag.String("pgpassword", "", "PostgreSQL password")
    pgDBName := flag.String("pgdbname", "", "PostgreSQL database name")
    flag.Parse()
    
    // Validate required parameters
    if *sqlitePath == "" || *pgUser == "" || *pgDBName == "" {
        fmt.Println("Usage: migrate -sqlite=path/to/sqlite.db -pguser=user -pgdbname=dbname [options]")
        flag.PrintDefaults()
        os.Exit(1)
    }
    
    // Connect to SQLite database
    sqliteConfig := database.NewSQLiteConfig(*sqlitePath)
    sqliteDriver, err := database.NewDriver(sqliteConfig)
    if err != nil {
        log.Fatalf("Failed to connect to SQLite database: %v", err)
    }
    
    // Connect to PostgreSQL database
    pgConfig := database.NewPostgresConfig(*pgHost, *pgPort, *pgUser, *pgPassword, *pgDBName, "disable")
    pgDriver, err := database.NewDriver(pgConfig)
    if err != nil {
        log.Fatalf("Failed to connect to PostgreSQL database: %v", err)
    }
    
    // Initialize PostgreSQL schema
    schemaManager := database.NewSchemaManager(pgDriver)
    if err := schemaManager.InitializeSchema(); err != nil {
        log.Fatalf("Failed to initialize PostgreSQL schema: %v", err)
    }
    
    // Migrate data
    if err := migrateData(sqliteDriver, pgDriver); err != nil {
        log.Fatalf("Migration failed: %v", err)
    }
    
    log.Println("Migration completed successfully")
}

func migrateData(source, target database.DBDriver) error {
    // Implement data migration logic for each table
    // ...
    return nil
}
```

## Next Steps

After implementing the database abstraction layer:

1. Test thoroughly with both database systems
2. Create PostgreSQL-specific schema optimizations
3. Implement data migration scripts
4. Set up monitoring and validation tools
5. Document the abstraction layer for future developers

## Conclusion

This implementation guide provides a foundation for creating a flexible database abstraction layer that will allow the application to work with both SQLite and PostgreSQL. By following this approach, you can gradually transition from SQLite to PostgreSQL while maintaining the ability to fall back to SQLite if needed.
