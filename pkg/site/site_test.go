package site

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	siteDir := filepath.Join(tempDir, "r2pm")

	s, err := New(siteDir)
	if err != nil {
		t.Fatal(err)
	}

	expectedDirs := []string{
		s.databaseSubDir(),
		s.gitSubDir(),
		s.installedSubDir(),
	}

	for _, ed := range expectedDirs {
		t.Run(ed, func(t *testing.T) {
			if info, err := os.Stat(ed); err != nil || !info.IsDir() {
				t.Fail()
			}
		})
	}
}
