package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var Database *DatabaseConfig
var Server *ServerConfig
var Auth *AuthConfig

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

type ServerConfig struct {
	Network  string
	Host     string
	Port     int
	Protocol string
	CertFile string
	KeyFile  string
}

type AuthConfig struct {
	SecretKey     string
	TokenDuration time.Duration
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func Load() {
	Server = &ServerConfig{
		Network:  getEnvAsString("SERVER_NETWORK", "tcp"),
		Host:     getEnvAsString("SERVER_HOST", "0.0.0.0"),
		Port:     getEnvAsInt("SERVER_PORT", 50051),
		Protocol: getEnvAsString("SERVER_PROTOCOL", "http"),
		CertFile: getEnvAsString("SERVER_CERT_FILE", ""),
		KeyFile:  getEnvAsString("SERVER_KEY_FILE", ""),
	}
	
	Database = &DatabaseConfig{
		Driver:   getEnvAsString("DATABASE_DRIVER", "mysql"),
		Host:     getEnvAsString("DATABASE_HOST", "localhost"),
		Port:     getEnvAsInt("DATABASE_PORT", 3306),
		Name:     getEnvAsString("DATABASE_NAME", "mydb"),
		User:     getEnvAsString("DATABASE_USER", "root"),
		Password: getEnvAsString("DATABASE_PASSWORD", "secret"),
	}
	
	Auth = &AuthConfig{
		SecretKey:     getEnvAsString("AUTH_SECRET_KEY", "secret"),
		TokenDuration: time.Duration(getEnvAsInt("AUTH_TOKEN_DURATION", 900)),
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
