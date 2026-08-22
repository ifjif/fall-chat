package utils

import "os"

func GetEnvWithDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}

func GetEnv(key string) (string, bool) {
	v, ok := os.LookupEnv(key)
	return v, ok
}
