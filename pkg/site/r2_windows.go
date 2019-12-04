package site

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (s Site) InstallRadare2(prefix, version string) error {
	url := fmt.Sprintf("http://radare.mikelloc.com/get/%s/radare2-msvc_64-%s.zip", version, version)

	fd, err := ioutil.TempFile("", "r2pm_*.zip")
	if err != nil {
		return fmt.Errorf("could not create a temporary file: %v", err)
	}
	defer fd.Close()

	log.Printf("Downloading %s into %s ", url, fd.Name())

	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error while downloading: %v", err)
	}
	defer res.Body.Close()

	n, err := bufio.NewReader(res.Body).WriteTo(fd)
	if err != nil {
		return fmt.Errorf("could not write the response body: %v", err)
	}

	fd.Seek(0, 0)

	z, err := zip.NewReader(fd, n)
	if err != nil {
		return fmt.Errorf("could not create a zip reader: %v", err)
	}

	const dirPerm = 0755

	log.Print("Extracting " + fd.Name())

	for _, f := range z.File {
		// Remove the first component of the path
		components := strings.SplitN(f.Name, "/", 2)
		if len(components) == 1 {
			// top-level directory - do nothing
			continue
		}

		target := filepath.Join(prefix, components[1])

		// Directory
		if f.FileInfo().IsDir() {
			log.Print("Creating " + target)

			if err := os.MkdirAll(target, dirPerm); err != nil {
				return err
			}

			continue
		}

		// File
		log.Print("Processing file " + target)

		dir := filepath.Dir(target)
		if err := os.MkdirAll(dir, dirPerm); err != nil {
			return err
		}

		zipFd, err := f.Open()
		if err != nil {
			return fmt.Errorf("could not extract %s: %v", f.Name, err)
		}

		fsFd, err := os.Create(target)
		if err != nil {
			return fmt.Errorf("could not create %s: %v", target, err)
		}

		if _, err := bufio.NewReader(zipFd).WriteTo(fsFd); err != nil {
			return fmt.Errorf("could not write %s: %v", target, err)
		}

		zipFd.Close()
		fsFd.Close()
	}

	return nil
}
