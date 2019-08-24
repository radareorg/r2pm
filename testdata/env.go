package testdata

import (
	"os"
	"testing"
)

func SetEnvVar(t *testing.T, name, value string) {
	t.Helper()

	if err := os.Setenv(name, value); err != nil {
		t.Fatalf("could not set %s: %v", name, err)
	}
}
