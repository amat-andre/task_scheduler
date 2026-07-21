package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server *ServerConfig
	DB     *DBConfig
	JWT    *JWTConfig
}

type ServerConfig struct {
	Port     string
	Password string
}

type DBConfig struct {
	Path string
}

type JWTConfig struct {
	Secret string
}

var defConfig = map[string]string{
	"TODO_PORT":     "7540",
	"TODO_DBFILE":   "scheduler.db",
	"TODO_PASSWORD": "",
	"JWT_SECRET":    "super_very_secret",
}

func Load() (*Config, error) {
	port, err := strconv.Atoi(getEnv("TODO_PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid TODO_PORT: %w", err)
	}

	return &Config{
		Server: &ServerConfig{
			Port:     fmt.Sprintf(":%d", port),
			Password: getEnv("TODO_PASSWORD"),
		},
		DB: &DBConfig{
			Path: getEnv("TODO_DBFILE"),
		},
		JWT: &JWTConfig{
			Secret: getEnv("JWT_SECRET"),
		},
	}, nil
}

func getEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defConfig[key]
}
