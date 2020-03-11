package r2package

import (
	"errors"
	"fmt"
	"io"
	"regexp"

	"gopkg.in/yaml.v2"
)

const (
	nameRE    = `^[a-zA-Z\-]+$`
	versionRE = `^\d+\.\d+\.\d+$`
)

var (
	namePattern    = regexp.MustCompile(nameRE)
	versionPattern = regexp.MustCompile(versionRE)
)

type ManagedFile struct {
	Path string
	Type string
}

type InstallInstructions struct {
	Out    []ManagedFile
	Source map[string]string
}

func (i *InstallInstructions) Check() error {
	if i.Source == nil || len(i.Source) == 0 {
		return errors.New("the source cannot be undefined")
	}

	return errors.New("not implemented")
}

type Package struct {
	Commands    []string
	Description string
	Install     map[string]InstallInstructions
	Name        string
	Tags        []string
	Version     string
}

func (p *Package) Check() error {
	if !namePattern.MatchString(p.Name) {
		return fmt.Errorf("name must match %s", nameRE)
	}

	if p.Description == "" {
		return errors.New("description cannot be empty")
	}

	if !versionPattern.MatchString(p.Version) {
		return fmt.Errorf("version should match the following regular expression: %s", versionRE)
	}

	if p.Install == nil || len(p.Install) == 0 {
		return errors.New("the manifest should contain the installation instructions for at least one operating system")
	}

	//for _, os := range p.Install {
	//
	//}

	return nil
}

func FromYAML(r io.Reader) (*Package, error) {
	p := &Package{}

	if err := yaml.NewDecoder(r).Decode(p); err != nil {
		return nil, fmt.Errorf("could not parse as a package: %v", err)
	}

	if err := p.Check(); err != nil {
		return nil, fmt.Errorf("invalid package description: %v", err)
	}

	return p, nil
}
