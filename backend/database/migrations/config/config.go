package configuration

type DatabaseConfig struct {
	Current struct {
		Driver string
		Path string   // For SQLite
		Host string	  // For PostgreSQL
		Port int
		User string
		Password string
		DBName string 
		SSLMode string
	}

	Target struct {
		Driver string 
		Host  string 
		Port int
		User string
		Password string
		DBName string
		SSLMode string
	}

	Migration struct {
		BatchSize int 
		TimeoutSeconds int
		RetryAttempts int
	}
}

func LoadConfig(configPath string) (*DatabaseConfig, error) {
	// TODO:: Implementation to load configuration from file
	return nil, nil
}