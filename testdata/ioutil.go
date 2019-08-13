package testdata

import (
	"io/ioutil"
	"testing"
)

func TempDirOrFail(t *testing.T) string {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	return dir
}
