package pkg

import "os"

func GetenvDefault(varName, defaultValue string) string {
	val := os.Getenv(varName)

	if val != "" {
		return val
	}

	return defaultValue
}
