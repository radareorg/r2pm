package dir

import (
	"testing"

	"github.com/radareorg/r2pm/testdata"
)

func testSiteDir(t *testing.T) {
	t.Run(`APPDATA=C:\temp`, func(t *testing.T) {
		testdata.SetEnvVar(t, "APPDATA", "/tmp/test")

		if s := SiteDir(); s != `C:\temp\RadareOrg\r2pm` {
			t.Fatal(s)
		}
	})
}
