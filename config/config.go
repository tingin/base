package config

import "os"

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func SetEnv(key, value string) error {
	return os.Setenv(key, value)
}
