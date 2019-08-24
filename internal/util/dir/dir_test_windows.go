package dir

import "testing"

func TestSiteDir(t *testing.T) {
	t.Run(`APPDATA=C:\temp`, func(t *testing.T) {
		if err := os.Setenv("APPDATA", `C:\temp`); err != nil {
			t.Fatalf("error while setting APPDATA: %v", err)
		}

		if s := SiteDir(); s != `C:\temp\RadareOrg\r2pm` {
			t.Fatal(s)
		}
	})
}
