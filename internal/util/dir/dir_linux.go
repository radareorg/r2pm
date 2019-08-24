package dir

import (
	"os"
	"path/filepath"
)

func platformPrefix() string {
	// Reference: https://specifications.freedesktop.org/basedir-spec/basedir-spec-0.7.html
	if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
		return xdgDataHome
	}

	return filepath.Join(os.Getenv("HOME"), ".local", "share")

}
