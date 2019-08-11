package pkg

import (
	"bufio"
	"os"
)

func CopyFile(src, dst string) error {
	srcFd, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFd.Close()

	dstFd, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFd.Close()

	_, err = bufio.NewReader(srcFd).WriteTo(dstFd)

	return err
}

func GetenvDefault(varName, defaultValue string) string {
	val := os.Getenv(varName)

	if val != "" {
		return val
	}

	return defaultValue
}
