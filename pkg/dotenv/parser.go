package dotenv

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadFromFile(filePath string) (map[string]string, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("environment file does not exist: %s", filePath)
	}

	envVars, err := godotenv.Read(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse environment file: %w", err)
	}

	return envVars, nil
}