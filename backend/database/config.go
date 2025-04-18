package database

// DBConfig holds database configuration
type DBConfig struct {
	// Common settings
	Driver string // "sqlite" or "postgres"
	DBName string

	// SQLite specific
	SQLitePath string
	
	// PostgreSQL specific
	Host string
	Port int
	UserName string
	Password string
	SSLMode string 

	// Connection pool settings
	MaxOpenConns int
	MaxIdleConns int
}

