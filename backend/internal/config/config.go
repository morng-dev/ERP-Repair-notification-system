package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	APPENV  string
	APPPORT string
	APPURL  string

	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMode string

	RedisHost string
	RedisPort string
	RedisPass string
	RedisDB   int

	JWTSecret    string
	JWTExpiresIn string
}

func LoadCongig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found, relying on environment variables")
	}
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	config := &Config{
		APPENV:    os.Getenv("APP_ENV"),
		APPPORT:   os.Getenv("APP_PORT"),
		APPURL:    os.Getenv("APP_URL"),
		DBHost:    os.Getenv("DB_HOST"),
		DBPort:    os.Getenv("DB_PORT"),
		DBUser:    os.Getenv("DB_USER"),
		DBPass:    os.Getenv("DB_PASS"),
		DBName:    os.Getenv("DB_NAME"),
		DBSSLMode: os.Getenv("DB_SSL"),

		RedisHost: os.Getenv("REDIS_HOST"),
		RedisPort: os.Getenv("REDIS_PORT"),
		RedisPass: os.Getenv("REDIS_PASS"),
		RedisDB:   db,

		JWTSecret:    os.Getenv("JWT_SECRET"),
		JWTExpiresIn: os.Getenv("JWT_EXPIRES_IN"),
	}
	return config
}
