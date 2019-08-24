package dir

import "testing"

func TestSiteDir(t *testing.T) {
	t.Run(`APPDATA=C:\temp`, func(t *testing.T) {
		testdata.SetEnvVar(t, "APPDATA", "/tmp/test")

		if s := SiteDir(); s != `C:\temp\RadareOrg\r2pm` {
			t.Fatal(s)
		}
	})
}
