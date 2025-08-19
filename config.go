package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	NWSBaseURL     string
	ServerPort     string
	RequestTimeout time.Duration
	ColdMax        int
	ModerateMax    int
}

func LoadConfig() (*Config, error) {
	envVars, err := loadEnvFile("conf.env")
	if err != nil {
		return nil, fmt.Errorf("failed to load conf.env: %w", err)
	}

	// Use sensible defaults if config missing
	config := &Config{
		NWSBaseURL:     getConfigValue(envVars, "NWS_BASE_URL", "https://api.weather.gov/"),
		ServerPort:     getConfigValue(envVars, "SERVER_PORT", "8000"),
		RequestTimeout: getDurationValue(envVars, "REQUEST_TIMEOUT", 30*time.Second),
		ColdMax:        getIntValue(envVars, "COLD_MAX", 50),
		ModerateMax:    getIntValue(envVars, "MODERATE_MAX", 80),
	}

	return config, nil
}

func loadEnvFile(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}
	defer file.Close()

	envVars := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			envVars[key] = value
		}
	}

	return envVars, scanner.Err()
}

func getConfigValue(envVars map[string]string, key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if value, exists := envVars[key]; exists {
		return value
	}
	return defaultValue
}

func getDurationValue(envVars map[string]string, key string, defaultValue time.Duration) time.Duration {
	valueStr := getConfigValue(envVars, key, "")
	if valueStr != "" {
		if duration, err := time.ParseDuration(valueStr); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getIntValue(envVars map[string]string, key string, defaultValue int) int {
	valueStr := getConfigValue(envVars, key, "")
	if valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
