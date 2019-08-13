package site

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/radareorg/r2pm/testdata"
)

func TestNew(t *testing.T) {
	tempDir := testdata.TempDirOrFail(t)
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
