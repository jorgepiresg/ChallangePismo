package config

import (
	"os"
	"strconv"
)

func New() Config {

	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))

	cfg := Config{
		ServerPort: os.Getenv("PORT"),
		DB: DB{
			MigrationFile: os.Getenv("DB_MIGRATION_FILE"),
			DriverName:    os.Getenv("DB_DRIVER_NAME"),
			Host:          os.Getenv("DB_HOST"),
			Port:          dbPort,
			User:          os.Getenv("POSTGRES_USER"),
			Password:      os.Getenv("POSTGRES_PASSWORD"),
			Database:      os.Getenv("POSTGRES_DB"),
		},
		Cache: Cache{
			Addr: os.Getenv("REDIS_ADDR"),
		},
	}

	return cfg
}

type Config struct {
	ServerPort string `json:"port"`
	DB         DB     `json:"db"`
	Cache      Cache  `json:"cache"`
}

type DB struct {
	MigrationFile string `json:"migration_file"`
	DriverName    string `json:"driver_name"`
	Host          string `json:"host"`
	Port          int    `json:"port"`
	User          string `json:"user"`
	Password      string `json:"password"`
	Database      string `json:"database"`
}

type Cache struct {
	Addr string `json:"addr"`
}
