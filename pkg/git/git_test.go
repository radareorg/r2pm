package git

import (
	"os"
	"strings"
	"testing"

	"github.com/radareorg/r2pm/pkg/process"
	"github.com/radareorg/r2pm/testdata"
)

func TestInit(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		tempDir := testdata.TempDirOrFail(t)
		defer os.RemoveAll(tempDir)

		if _, err := Init(tempDir, false); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("recreate, force=false", func(t *testing.T) {
		tempDir := testdata.TempDirOrFail(t)
		defer os.RemoveAll(tempDir)

		if _, err := Init(tempDir, false); err != nil {
			t.Fatal(err)
		}

		if _, err := Init(tempDir, false); err == nil {
			testdata.FailExpectedError(t)
		}
	})

	t.Run("recreate, force=true", func(t *testing.T) {
		tempDir := testdata.TempDirOrFail(t)
		defer os.RemoveAll(tempDir)

		if _, err := Init(tempDir, false); err != nil {
			t.Fatal(err)
		}

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

func TestRepository_AddRemote(t *testing.T) {
	const (
		name = "random-name"
		url  = "git@localhost.localdomain:test"
	)

	tempDir := testdata.TempDirOrFail(t)
	defer os.RemoveAll(tempDir)

	repo, err := Init(tempDir, false)
	if err != nil {
		t.Fatal(err)
	}

	if err := repo.AddRemote(name, url); err != nil {
		t.Fatal(err)
	}

	res, err := process.Run("git", []string{"remote", "-v"}, tempDir)
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(res.Stdout.String(), "\n")

	// maybe not very portable...
	for i, l := range lines {
		if l == "" {
			continue
		}

		// remotes are printed in the following way: "<name>\t<url> <(fetch|push)>
		beforeSpace := strings.Split(l, " ")[0]

		words := strings.Split(beforeSpace, "\t")
		w := words[0]
		if w != name {
			t.Fatalf("%q: word %d:0 is %q, not %q", beforeSpace, i+1, w, name)
		}

		w = words[1]
		if w != url {
			t.Fatalf("%q: word %d:1 is %q, not %q", beforeSpace, i+1, w, url)
		}
	}

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
