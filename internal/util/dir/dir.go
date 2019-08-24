package dir

import (
	"os"
	"path/filepath"
)

const SiteDirEnvVar = "R2PM_SITEDIR"

func SiteDir() string {
	if envVar := os.Getenv(SiteDirEnvVar); envVar != "" {
		return envVar
	}

	return filepath.Join(platformPrefix(), "RadareOrg", "r2pm")
}
