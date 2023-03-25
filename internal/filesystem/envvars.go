package filesystem

import (
	"errors"
	"fmt"
	"os"
)

type EnvVars map[string]string

func CheckEnvVars(keys []string) error {
	for _, key := range keys {
		if _, ok := os.LookupEnv(key); !ok {
			return errors.New(fmt.Sprintf("Environment variable %s is not set", key))
		}
	}
	return nil
}

func IsEnvVarSetOrExported(key string) error {
	if _, ok := os.LookupEnv(key); !ok {
		return errors.New(fmt.Sprintf("Environment variable %s is not set", key))
	}
	return nil
}

func GetEnvVars(keys []string) EnvVars {
	result := make(EnvVars)
	for _, key := range keys {
		if value, ok := os.LookupEnv(key); ok {
			result[key] = value
		}
	}
	return result
}

func ScanAWSCredsFromEnv() EnvVars {
	keys := []string{
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
	}
	return GetEnvVars(keys)
}
