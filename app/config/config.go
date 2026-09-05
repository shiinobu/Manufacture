package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Load() error {
	if err := godotenv.Load(); err != nil {
		// .env is optional when environment variables are provided by the runtime.
	}

	if os.Getenv("JWT_KEY") == "" {
		return fmt.Errorf("JWT_KEY environment variable is required")
	}

	return nil
}

func JWTKey() []byte {
	return []byte(os.Getenv("JWT_KEY"))
}
