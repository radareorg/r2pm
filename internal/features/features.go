package features

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/radareorg/r2pm/pkg/r2package"
	"github.com/radareorg/r2pm/pkg/site"
)

const (
	DebugEnvVar = "R2PM_DEBUG"

	msgCannotInitialize = "could not initialize: %w"
)

func Delete(r2pmDir string) error {
	s, err := site.New(r2pmDir)
	if err != nil {
		return fmt.Errorf(msgCannotInitialize, err)
	}

	return s.Remove()
}

func Init(r2pmDir string) error {
	s, err := site.New(r2pmDir)
	if err != nil {
		return fmt.Errorf("could not initialize: %w", err)
	}

	return s.Database().InitOrUpdate()
}

func Install(r2pmDir, packageName string) error {
	s, err := site.New(r2pmDir)
	if err != nil {
		return fmt.Errorf(msgCannotInitialize, err)
	}

	return s.InstallPackage(packageName)
}

func InstallFromFile(r2pmDir, path string) error {
	s, err := site.New(r2pmDir)
	if err != nil {
		return fmt.Errorf(msgCannotInitialize, err)
	}

	return s.InstallPackageFromFile(path)
}

func InstallRadare2(r2pmDir, r2Dir, version string) error {
	s, err := site.New(r2pmDir)
	if err != nil {
		return err
	}

	return s.InstallRadare2(r2Dir, version)
}

func UninstallRadare2(r2pmDir, r2Dir string) error {
	s, err := site.New(r2pmDir)
	if err != nil {
		return err
	}

	return s.UninstallRadare2(r2Dir)
}

func ListAvailable(r2pmDir string) ([]r2package.Info, error) {
	s, err := site.New(r2pmDir)
	if err != nil {
		return nil, fmt.Errorf(msgCannotInitialize, err)
	}

	return s.Database().ListAvailablePackages()
}

func ListInstalled(r2pmDir string) ([]r2package.Info, error) {
	s, err := site.New(r2pmDir)
	if err != nil {
		return nil, fmt.Errorf(msgCannotInitialize, err)
	}

	return s.ListInstalledPackages()
}

func Search(r2pmDir, pattern string) ([]r2package.Info, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("%q is not a valid regex: %w", pattern, err)
	}

	packages, err := ListAvailable(r2pmDir)
	if err != nil {
		return nil, fmt.Errorf("could not get the list of packages: %w", err)
	}

	matches := make([]r2package.Info, 0, len(packages))

	for _, p := range packages {
		if re.Match([]byte(p.Name)) {
			matches = append(matches, p)
		}
	}

	return matches, nil
}

func SetDebug(value bool) {
	if value {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(ioutil.Discard)
	}
}

func Uninstall(r2pmDir, packageName string) error {
	s, err := site.New(r2pmDir)
	if err != nil {
		return fmt.Errorf(msgCannotInitialize, err)
	}

	return s.UninstallPackage(packageName)
}

func Upgrade(r2pmDir, packageName string) error {
	s, err := site.New(r2pmDir)
	if err != nil {
		return fmt.Errorf(msgCannotInitialize, err)
	}

	return s.Upgrade(packageName)
}

func UpgradeAll(r2pmDir string) error {
	s, err := site.New(r2pmDir)
	if err != nil {
		return fmt.Errorf(msgCannotInitialize, err)
	}

	packages, err := s.ListInstalledPackages()
	if err != nil {
		log.Print(err)
		return errors.New("could not list the installed packages")
	}

	failed := make([]string, 0, len(packages))

	for _, p := range packages {
		name := p.Name

		log.Println("Upgrading " + name)

		if err := s.Upgrade(name); err != nil {
			log.Print(err)
			failed = append(failed, name)
		}
	}

	sort.Strings(failed)

	if len(failed) > 0 {
		return fmt.Errorf(
			"could not upgrade the following packages: %s",
			strings.Join(failed, ", "))
	}

	return nil
}
