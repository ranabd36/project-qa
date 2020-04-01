package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Database DatabaseConfig
}

type DatabaseConfig struct {
	DatabaseDriver     string
	DatabaseHost     string
	DatabasePort     int
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func Get() *Config {
	return &Config{
		Database: DatabaseConfig{
			DatabaseDriver:     getEnvAsString("DATABASE_DRIVER", "mysql"),
			DatabaseHost:     getEnvAsString("DATABASE_HOST", "localhost"),
			DatabasePort:     getEnvAsInt("DATABASE_PORT", 3306),
			DatabaseName:     getEnvAsString("DATABASE_NAME", "mydb"),
			DatabaseUser:     getEnvAsString("DATABASE_USER", "root"),
			DatabasePassword: getEnvAsString("DATABASE_PASSWORD", "secret"),
		},
	}
}

func getEnvAsString(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(name string, defaultValue int) int {
	valueStr := getEnvAsString(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	
	return defaultValue
}

func getEnvAsBool(name string, defaultValue bool) bool {
	valStr := getEnvAsString(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	
	return defaultValue
}

func getEnvAsSlice(name string, defaultValue []string, sep string) []string {
	valStr := getEnvAsString(name, "")
	
	if valStr == "" {
		return defaultValue
	}
	
	val := strings.Split(valStr, sep)
	
	return val
}
