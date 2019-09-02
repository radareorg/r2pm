// +build integration

package site

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestInstallRadare2(t *testing.T) {
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

	t.Log("Installing")

	if err := s.InstallRadare2(prefix); err != nil {
		t.Fatal(err)
	}

	r2Path := filepath.Join(prefix, "bin", "r2")

	if _, err := os.Stat(r2Path); err != nil {
		t.Fatalf("Could not stat(%q)", r2Path)
	}
}
