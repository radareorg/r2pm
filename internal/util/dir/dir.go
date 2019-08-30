package dir

import (
	"os"
	"path/filepath"
)

const (
	orgSubDir     = "RadareOrg"
	SiteDirEnvVar = "R2PM_SITEDIR"
)

func R2Dir() string {
	return filepath.Join(orgSubdDir(), "radare2")
}

func SiteDir() string {
	if envVar := os.Getenv(SiteDirEnvVar); envVar != "" {
		return envVar
	}

	return filepath.Join(orgSubdDir(), "r2pm")
}

func orgSubdDir() string {
	return filepath.Join(platformPrefix(), orgSubDir)
}
