package config

import (
	"fmt"
	"os"
)

// Config содержит конфигурацию приложения
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

// Load загружает конфигурацию из .env файла
func Load() (*Config, error) {
	return LoadConfig()
}

func LoadConfig() (*Config, error) {
	config := &Config{
		DBHost:     getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:     getEnvOrDefault("DB_PORT", "5432"),
		DBUser:     getEnvOrDefault("DB_USER", "postgres"),
		DBPassword: getEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName:     getEnvOrDefault("DB_NAME", "music_library"),
		ServerPort: getEnvOrDefault("SERVER_PORT", "8080"),
	}

	return config, nil
}

// GetDBConnString возвращает строку подключения к базе данных
func (c *Config) GetDBConnString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

// GetDBConnStringWithoutDatabase возвращает строку подключения к postgres без указания базы данных
func (c *Config) GetDBConnStringWithoutDatabase() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword)
}

// getEnvOrDefault возвращает значение переменной окружения или значение по умолчанию
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
