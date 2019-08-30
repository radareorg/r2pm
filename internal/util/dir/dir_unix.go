// +build darwin freebsd

package dir

func platformPrefix() string {
	return "/usr/local/share"
}
