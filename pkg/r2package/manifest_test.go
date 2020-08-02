package r2package

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/radareorg/r2pm/pkg/r2package/fetchers"
)

func TestManifest_Verify(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		m := &Manifest{}

		if err := m.Verify(); err == nil {
			t.Fatal("Expected an error")
		}
	})

	t.Run("empty version", func(t *testing.T) {
		m := Manifest{Name: "test-package"}

		if err := m.Verify(); err == nil {
			t.Fatal("Expected an error")
		}
	})

	t.Run("empty description", func(t *testing.T) {
		m := Manifest{
			Name:    "test-package",
			Version: "1.2.3",
		}

		if err := m.Verify(); err == nil {
			t.Fatal("Expected an error")
		}
	})

	t.Run("empty install instructions", func(t *testing.T) {
		m := Manifest{
			Name:        "test-package",
			Version:     "1.2.3",
			Description: "Test package",
		}

		if err := m.Verify(); err == nil {
			t.Fatal("Expected an error")
		}
	})
}

func Test_getFetcher(t *testing.T) {
	t.Run("git", func(t *testing.T) {
		const (
			ref  = "master"
			repo = "git@remote.tld:path/to/repo.git"
		)

		s := Source{
			Type: "git",
			Ref:  ref,
			Repo: repo,
		}

		f, err := getFetcher(s)
		if err != nil {
			t.Fatalf("Could not get the fetcher: %v", err)
		}

		if _, ok := f.(*fetchers.Git); !ok {
			t.Fatalf("Unexpected fetcher type %T", f)
		}
	})

	t.Run("zip", func(t *testing.T) {
		const url = "https://domain.tld/some-archive.zip"

		s := Source{
			Type: "zip",
			URL:  url,
		}

		f, err := getFetcher(s)
		if err != nil {
			t.Fatalf("Could not get the fetcher: %v", err)
		}

		if _, ok := f.(*fetchers.Zip); !ok {
			t.Fatalf("Unexpected fetcher type %T", f)
		}
	})

	t.Run("not-a-fetcher", func(t *testing.T) {
		s := Source{Type: "not-a-source"}

		if _, err := getFetcher(s); err == nil {
			t.Fatal("Expected an error")
		}
	})
}

func TestFromYAML(t *testing.T) {
	r := strings.NewReader(`---
name: test-package
version: 1.2.3
description: test package
install:
  linux:
    source:
      type: git
      repo: git@remote.tld:path/to/repo.git
      ref: master
    commands:
      - linux command
    out:
      - path: relative/path/to/out/file/linux
        type: shared-lib
  windows:
    source:
      type: zip
      url: https://domain.tld/some-archive.zip
    commands:
      - windows command
    out:
      - path: relative/path/to/out/file/windows
        type: shared-lib
`)

	got, err := FromYAML(r)
	if err != nil {
		t.Fatalf("Could not parse the input: %v", err)
	}

	want := &Manifest{
		Name:        "test-package",
		Version:     "1.2.3",
		Description: "test package",
		Install: map[string]InstallInstructions{
			"linux": {
				Commands: []string{"linux command"},
				Out: []ManagedFile{
					{
						Path: "relative/path/to/out/file/linux",
						Type: "shared-lib",
					},
				},
				Source: Source{
					Type: "git",
					Ref:  "master",
					Repo: "git@remote.tld:path/to/repo.git",
				},
			},
			"windows": {
				Commands: []string{"windows command"},
				Out: []ManagedFile{
					{
						Path: "relative/path/to/out/file/windows",
						Type: "shared-lib",
					},
				},
				Source: Source{
					Type: "zip",
					URL:  "https://domain.tld/some-archive.zip",
				},
			},
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatal(diff)
	}
}

func Test_osFromGOOS(t *testing.T) {
	cases := []struct {
		goos, expected string
		returnsError   bool
	}{
		{
			goos:     "darwin",
			expected: "macos",
		},
		{
			goos:     "freebsd",
			expected: "freebsd",
		},
		{
			goos:     "linux",
			expected: "linux",
		},
		{
			goos:         "openbsd",
			returnsError: true,
		},
		{
			goos:     "windows",
			expected: "windows",
		},
	}

	for _, c := range cases {
		t.Run(c.goos, func(t *testing.T) {
			got, err := osFromGOOS(c.goos)
			if err != nil {
				if c.returnsError {
					return
				}

				t.Fatalf("Unexpected error: %got", err)
			}

			if got != c.expected {
				t.Fatalf("Expected %q, got %q", c.expected, got)
			}
		})
	}
}
