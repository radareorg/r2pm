package features

import (
	"log"
	"regexp"

	"golang.org/x/xerrors"

	"github.com/radareorg/r2pm/pkg/database"
)

func Delete(r2pmDir string) error {
	return database.Delete(r2pmDir)
}

func Init(r2pmDir string) error {
	return database.Init(r2pmDir)
}

func Install(r2pmDir, packageName string) error {
	pi, err := database.FindPackage(r2pmDir, packageName)
	if err != nil {
		log.Fatalf("could not find package %s: %v", packageName, err)
	}

	if err := pi.Install(r2pmDir); err != nil {
		log.Fatalf("could not install %s: %v", packageName, err)
	}

	return nil
}

// TODO this should return a slice of r2package.Info
func List(r2pmDir string) ([]string, error) {
	return database.List(r2pmDir)
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
	pi, err := database.FindPackage(r2pmDir, packageName)
	if err != nil {
		log.Fatalf("could not find package %s: %v", packageName, err)
	}

	if err := pi.Uninstall(r2pmDir); err != nil {
		log.Fatalf("could not uninstall %s: %v", packageName, err)
	}

	return nil
}
