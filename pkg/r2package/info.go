package r2package

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

//
// Info
//

type OutputFile struct {
	Path	string
	Type	string
}

type Info struct {
	Name		string
	Version		string
	Description	string
	Tags		[]string
	// avoid conflict with i.Install()
	InstallConf struct {
		Linux struct {
			Source struct {
				Type	string
				Url	string	// for zip
				Repo	string	// for git
				Ref	string
			}
			Commands []string
			Out []OutputFile `yaml:"out,flow"`	// for zip
		}
		Macos struct {
			Source struct {
				Type	string
				Url	string
				Repo	string
				Ref	string
			}
			Commands []string
			Out []OutputFile `yaml:"out,flow"`
		}
		Windows struct {
			Source struct {
				Type	string
				Url	string
				Repo	string
				Ref	string
			}
			Commands []string
			Out []OutputFile `yaml:"out,flow"`
		}
	} `yaml:"install"`
	// TODO: windows, macos, uninstall, out, tags
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
	// TODO: don't hardcode Linux
	platform := i.InstallConf.Linux
	switch platform.Source.Type {
	case "git":
		return gitInstaller{i}, nil
	// TODO: zip
	default:
		return nil, fmt.Errorf("%q: unhandled package type",
					platform.Source.Type)
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
		return nil, fmt.Errorf("could not read %s: %w", path, err)
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
			fmt.Printf("Warning: could not read %s: %v", name, err)
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
