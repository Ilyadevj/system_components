package config

import "os"

type Config struct {
	DatabaseDSN string
}

func Load() Config {
	dsn := os.Getenv("DATABASE_DSN")

	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=system_components port=5432 sslmode=disable"
	}

	return Config{
		DatabaseDSN: dsn,
	}
}