package features

import (
	"regexp"

	"golang.org/x/xerrors"

	"github.com/radareorg/r2pm/pkg/site"
)

const msgCannotInitialize = "could not initialize: %w"

func Delete(r2pmDir string) error {
	s, err := site.New(r2pmDir)
	if err != nil {
		return xerrors.Errorf(msgCannotInitialize, err)
	}

	return s.Remove()
}

func Init(r2pmDir string) error {
	s, err := site.New(r2pmDir)
	if err != nil {
		return xerrors.Errorf("could not initialize: %w", err)
	}

	return s.Database().InitOrUpdate()
}

func Install(r2pmDir, packageName string) error {
	s, err := site.New(r2pmDir)
	if err != nil {
		return xerrors.Errorf(msgCannotInitialize, err)
	}

	return s.InstallPackage(packageName)
}

// TODO this should return a slice of r2package.Info
func List(r2pmDir string) ([]string, error) {
	s, err := site.New(r2pmDir)
	if err != nil {
		return nil, xerrors.Errorf(msgCannotInitialize, err)
	}

	return s.ListInstalledPackages()
}

// TODO this should return a slice of r2package.Info
func Search(r2pmDir, pattern string) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, xerrors.Errorf("%q is not a valid regex: %w", pattern, err)
	}

	packages, err := List(r2pmDir)
	if err != nil {
		return nil, xerrors.Errorf("could not get the list of packages: %w", err)
	}

	matches := make([]string, 0, len(packages))

	for _, p := range packages {
		if re.Match([]byte(p)) {
			matches = append(matches, p)
		}
	}

	return matches, nil
}

func Uninstall(r2pmDir, packageName string) error {
	s, err := site.New(r2pmDir)
	if err != nil {
		return xerrors.Errorf(msgCannotInitialize, err)
	}

	return s.UninstallPackage(packageName)
}
