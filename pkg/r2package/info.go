package r2package

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
	"gopkg.in/yaml.v2"
)

//
// Info
//

type Info struct {
	Name          string
	Type          string
	Repo          string
	Desc          string
	InstallCmds   []string `yaml:"install"`
	UninstallCmds []string `yaml:"uninstall"`
}

func (i Info) Install(inDir string) error {
	installer, err := i.installer()
	if err != nil {
		return err
	}

	return installer.install(inDir)
}

func (i Info) Uninstall(inDir string) error {
	installer, err := i.installer()
	if err != nil {
		return err
	}

	return installer.uninstall(inDir)
}

func (i Info) installer() (installer, error) {
	switch i.Type {
	case "git":
		return gitInstaller{i}, nil
	default:
		return nil, xerrors.Errorf("%q: unhandled package type", i.Type)
	}
}

//
// InfoFile
//

type InfoFile struct {
	Info
	Path string
}

func FromFile(path string) (InfoFile, error) {
	infoFile := InfoFile{Path: path}

	fd, err := os.Open(path)
	if err != nil {
		return infoFile, err
	}
	defer fd.Close()

	d := yaml.NewDecoder(fd)
	d.SetStrict(true)

	err = d.Decode(&infoFile.Info)

	return infoFile, err
}

func ReadDir(path string) ([]InfoFile, error) {
	log.Println("Reading " + path)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, xerrors.Errorf("could not read %s: %w", path, err)
	}

	packages := make([]InfoFile, 0, len(files))

	for _, f := range files {
		// skip directories
		if f.IsDir() {
			continue
		}

		name := filepath.Join(path, f.Name())

		ifile, err := FromFile(name)
		if err != nil {
			log.Printf("could not read %s: %w", name, err)
			continue
		}

		packages = append(packages, ifile)
	}

	return packages, nil
}

//
// installer
//

type installer interface {
	install(string) error
	uninstall(string) error
}
