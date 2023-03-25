package filesystem

import (
	"fmt"
	"os"
)

func ReadYAMLFile(filepath string) ([]byte, error) {
	yamlContent, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	return yamlContent, nil
}
