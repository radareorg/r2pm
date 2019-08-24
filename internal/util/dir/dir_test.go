package dir

import (
	"os"
	"testing"
)

func setEnvVar(t *testing.T, name, value string) {
	t.Helper()

	if err := os.Setenv(name, value); err != nil {
		t.Fatalf("could not set %s: %v", name, err)
	}
}

func TestSiteDir(t *testing.T) {
	// this makes the test mandatory
	testSiteDir(t)
}
