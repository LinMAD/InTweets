package core

import (
	"fmt"
	"os"
)

// GetEnvVar value
func GetEnvVar(n string) (string, error) {
	v := os.Getenv(n)
	if v == "" {
		return "", fmt.Errorf("missing environment variable: " + n)
	}

	return v, nil
}
