package config

import (
	"os"
	"strconv"
	"strings"
)

type KafkaConfig struct {
	Enabled          bool
	Brokers          []string
	Topic            string
	GroupID          string
	SecurityProtocol string
	SASLMechanism    string
	SASLUsername     string
	SASLPassword     string
	TimeoutSeconds   int
	RetryMax         int
}

func NewKafkaConfig() KafkaConfig {
	return KafkaConfig{
		Enabled:          getBoolEnv("KAFKA_ENABLED", false),
		Brokers:          parseStringSliceEnv("KAFKA_BROKERS", []string{"localhost:9092"}),
		Topic:            getEnv("KAFKA_TOPIC", ""),
		GroupID:          getEnv("KAFKA_GROUP_ID", "erp-consumer"),
		SecurityProtocol: getEnv("KAFKA_SECURITY_PROTOCOL", "plaintext"),
		SASLMechanism:    getEnv("KAFKA_SASL_MECHANISM", "PLAIN"),
		SASLUsername:     getEnv("KAFKA_SASL_USERNAME", ""),
		SASLPassword:     getEnv("KAFKA_SASL_PASSWORD", ""),
		TimeoutSeconds:   getIntEnv("KAFKA_TIMEOUT_SECONDS", 30),
		RetryMax:         getIntEnv("KAFKA_RETRY_MAX", 3),
	}
}

func (c KafkaConfig) IsConfigured() bool {
	return c.Enabled && len(c.Brokers) > 0 && c.Topic != ""
}

func getEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func getBoolEnv(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getIntEnv(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func parseStringSliceEnv(key string, fallback []string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}

	if len(result) > 0 {
		return result
	}

	return fallback
}
