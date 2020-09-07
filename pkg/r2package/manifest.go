package r2package

import (
	"errors"
	"fmt"
	"io"

	"gopkg.in/yaml.v2"

	"github.com/radareorg/r2pm/pkg/r2package/fetchers"
)

type ManagedFile struct {
	Path string
	Type string
}

const (
	gitSource = "git"
	zipSource = "zip"
)

type Source struct {
	Type string

	// type: git
	Ref  string
	Repo string

	// type: zip
	URL string
}

type InstallInstructions struct {
	Commands []string
	Out      []ManagedFile
	Source   Source
}

type Manifest struct {
	Name        string
	Version     string
	Description string
	Install     map[string]InstallInstructions
}

func (m *Manifest) Verify() error {
	if m.Name == "" {
		return errors.New("name cannot be empty")
	}

	if m.Version == "" {
		return errors.New("version cannot be empty")
	}

	if m.Description == "" {
		return errors.New("description cannot be empty")
	}

	if len(m.Install) == 0 {
		return errors.New("no installation instructions")
	}

	return nil
}

func FromYAML(r io.Reader) (*Manifest, error) {
	m := &Manifest{}

	return m, yaml.NewDecoder(r).Decode(m)
}

func getFetcher(s Source) (fetcher, error) {
	switch s.Type {
	case gitSource:
		return fetchers.NewGit(s.Repo, s.Ref), nil
	case zipSource:
		return fetchers.NewZip(s.URL), nil
	default:
		return nil, fmt.Errorf("%q: invalid source", s.Type)
	}
}

func osFromGOOS(goos string) (string, error) {
	switch goos {
	case "darwin":
		return "macos", nil
	case "freebsd", "linux", "windows":
		return goos, nil
	default:
		return "", fmt.Errorf("%q: unhandled OS", goos)
	}
}
