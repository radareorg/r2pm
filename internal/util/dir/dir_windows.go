package dir

func platformPrefix() string {
	var prefix string

	if appData := os.Getenv("APPDATA"); appData != "" {
		return appData
	}

	return os.Getenv("HOMEPATH")
}
