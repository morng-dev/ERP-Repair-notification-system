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

	SMTP_HOST     string
	SMTP_PORT     int
	SMTP_USERNAME string
	SMTP_PASSWORD string
	SMTP_FROM     string
}

func LoadCongig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found, relying on environment variables")
	}
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Printf("Warning: REDIS_DB found, relying on environment variables")
	}
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Printf("Warning: SMTP_PORT found, relying on environment variables")
	}
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

		SMTP_HOST:     os.Getenv("SMTP_HOST"),
		SMTP_PORT:     smtpPort,
		SMTP_USERNAME: os.Getenv("SMTP_USERNAME"),
		SMTP_PASSWORD: os.Getenv("SMTP_PASSWORD"),
		SMTP_FROM:     os.Getenv("SMTP_FROM"),
	}
	return config
}
