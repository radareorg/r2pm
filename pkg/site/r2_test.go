// +build integration

package site

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestRadare2(t *testing.T) {
	siteDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(siteDir)

	s, err := New(siteDir)
	if err != nil {
		t.Fatal(err)
	}

	prefix, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(prefix)

	t.Run("InstallRadare2", func(t *testing.T) {
		version := os.Getenv("R2_VERSION")

		if version == "" {
			t.Fatal("The r2 version must be defined")
		}

		if err := s.InstallRadare2(prefix, version); err != nil {
			t.Fatal(err)
		}

		r2Bin := "r2"

		if runtime.GOOS == "windows" {
			r2Bin = "r2.bat"
		}

		r2Path := filepath.Join(prefix, "bin", r2Bin)

		if _, err := os.Stat(r2Path); err != nil {
			t.Fatalf("Could not stat(%q)", r2Path)
		}
	})

	t.Run("UninstallRadare2", func(t *testing.T) {
		if err := s.UninstallRadare2(prefix); err != nil {
			t.Fatal(err)
		}
	})
}
