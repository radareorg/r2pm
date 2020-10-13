package features

import (
	"strings"
)

func isInvalidDirectory(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "no such file or directory")
}
