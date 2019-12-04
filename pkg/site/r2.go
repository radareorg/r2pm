package site

import "os"

// UninstallRadare2 removes the directory containing the radare2 installation and
// all the files it contains.
func (s Site) UninstallRadare2(prefix string) error {
	return os.RemoveAll(prefix)
}
