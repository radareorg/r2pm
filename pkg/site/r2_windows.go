package site

import (
	"archive/zip"
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/xerrors"
)

func (s Site) InstallRadare2(prefix string) error {
	url := "http://radare.mikelloc.com/get/3.6.0/radare2-msvc_64-3.6.0.zip"

	fd, err := ioutil.TempFile("", "r2pm_*")
	if err != nil {
		return xerrors.Errorf("could not create a temporary file: %v", err)
	}
	defer fd.Close()

	log.Printf("Downloading %s into %s ", fd.Name(), url)

	res, err := http.Get(url)
	if err != nil {
		return xerrors.Errorf("error while downloading: %v", err)
	}
	defer res.Body.Close()

	n, err := bufio.NewReader(res.Body).WriteTo(fd)
	if err != nil {
		return xerrors.Errorf("could not write the response body: %v", err)
	}

	fd.Seek(0, 0)

	z, err := zip.NewReader(fd, n)
	if err != nil {
		return xerrors.Errorf("could not create a zip reader: %v", err)
	}

	const dirPerm = 0755

	for _, f := range z.File {
		// Remove the first component of the path
		components := strings.SplitN(f.Name, string(os.PathSeparator), 1)
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
		log.Print("Extracting " + target)

		dir := filepath.Dir(target)
		if err := os.MkdirAll(dir, dirPerm); err != nil {
			return err
		}

		zipFd, err := f.Open()
		if err != nil {
			return xerrors.Errorf("could not extract %s: %v", f.Name, err)
		}
		defer zipFd.Close()

		fsFd, err := os.Create(target)
		if err != nil {
			return xerrors.Errorf("could not create %s: %v", target, err)
		}
		defer fsFd.Close()

		if _, err := bufio.NewReader(zipFd).WriteTo(fsFd); err != nil {
			return xerrors.Errorf("could not write %s: %v", target, err)
		}
	}

	return nil
}
