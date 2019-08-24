package dir

import "testing"

func testSiteDir(t *testing.T) {
	if s := SiteDir(); s != "/usr/share/RadareOrg/r2pm" {
		t.Fatal(s)
	}
}
