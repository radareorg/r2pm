package r2package

import (
	"os"

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

//
// installer
//

type installer interface {
	install(string) error
	uninstall(string) error
}
