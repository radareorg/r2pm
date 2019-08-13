package git

import (
	"os"
	"testing"

	"github.com/radareorg/r2pm/testdata"
)

func TestInit(t *testing.T) {
	tempDir := testdata.TempDirOrFail(t)
	defer os.RemoveAll(tempDir)

	t.Run("create", func(t *testing.T) {
		if _, err := Init(tempDir, false); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("recreate, force=false", func(t *testing.T) {
		if _, err := Init(tempDir, false); err == nil {
			t.Fatal("should return an error")
		}
	})

	t.Run("recreate, force=true", func(t *testing.T) {
		if _, err := Init(tempDir, true); err != nil {
			t.Fatal(err)
		}
	})
}

func TestOpen(t *testing.T) {
	t.Run("/dev/null", func(t *testing.T) {
		if _, err := Open(os.DevNull); err == nil {
			testdata.FailExpectedError(t)
		}
	})

	t.Run("not a git repo", func(t *testing.T) {
		tempDir := testdata.TempDirOrFail(t)
		defer os.RemoveAll(tempDir)

		if _, err := Open(tempDir); err == nil {
			testdata.FailExpectedError(t)
		}
	})

	t.Run("should work", func(t *testing.T) {
		tempDir := testdata.TempDirOrFail(t)
		defer os.RemoveAll(tempDir)

		if _, err := Init(tempDir, false); err != nil {
			t.Fatalf("could not create the repository: %v", err)
		}

		if _, err := Open(tempDir); err != nil {
			t.Fatal(err)
		}
	})
}

func TestRun(t *testing.T) {
	tempDir := testdata.TempDirOrFail(t)
	defer os.RemoveAll(tempDir)

	// check that we cannot open the repository yet
	if _, err := Open(tempDir); err == nil {
		testdata.FailExpectedError(t)
	}

	if err := Run([]string{"init"}, tempDir); err != nil {
		t.Fatal(err)
	}

	// now we should be able to open the directory
	if _, err := Open(tempDir); err != nil {
		t.Fatal(err)
	}
}
