package dir

import "os"

func platformPrefix() string {
	if appData := os.Getenv("APPDATA"); appData != "" {
		return appData
	}

	return os.Getenv("HOMEPATH")
}
