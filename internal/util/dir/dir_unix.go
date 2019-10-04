// +build darwin freebsd openbsd

package dir

func platformPrefix() string {
	return "/usr/local/share"
}
