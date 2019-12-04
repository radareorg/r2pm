// +build darwin freebsd

package dir

import (
	"os"
	"path/filepath"
	"testing"
)

func testSiteDir(t *testing.T) {
	if s := SiteDir(); s != filepath.Join(os.Getenv("HOME"), "bin", "prefix", "RadareOrg", "r2pm") {
		t.Fatal(s)
	}
}
