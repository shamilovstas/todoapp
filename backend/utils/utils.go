package utils

import (
	"fmt"
	"os"
)

func GetEnvWithDefault(key, def string) string {
	envValue, err := GetEnv(key)
	if err != nil {
		envValue = def
	}
	return envValue
}

type ErrEnvMissing struct {
	key string
}

func (e ErrEnvMissing) Error() string {
	return fmt.Sprintf("env key %s is missing", e.key)
}

func GetEnv(key string) (string, error) {
	envValue := os.Getenv(key)
	if len(envValue) == 0 {
		return "", &ErrEnvMissing{key: key}
	}
	return envValue, nil
}
