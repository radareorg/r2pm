package features

import (
	"strings"
)

func isInvalidDirectory(err error) bool {
	return strings.Contains(err.Error(), "no such file or directory")
}
