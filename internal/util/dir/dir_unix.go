// +build darwin freebsd openbsd

package dir

import (
	"os"
	"path/filepath"
)

func platformPrefix() string {
	return filepath.Join(os.Getenv("HOME"), "bin", "prefix")
}
