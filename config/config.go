package config

import "os"

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	PGADMINDefaultEmail    string
	PGADMINDefaultPassword string

	RedisHost string
	RedisPort string
}

func LoadConfig() *Config {
	return &Config{
		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBUser: os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName: os.Getenv("DB_NAME"),

		PGADMINDefaultEmail: os.Getenv("PGADMIN_DEFAULT_EMAIL"),
		PGADMINDefaultPassword: os.Getenv("PGADMIN_DEFAULT_PASSWORD"),

		RedisHost: os.Getenv("REDIS_HOST"),
		RedisPort: os.Getenv("REDIS_PORT"),
	}
}