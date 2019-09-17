package site

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/radareorg/r2pm/pkg"
	"github.com/radareorg/r2pm/pkg/database"
	"github.com/radareorg/r2pm/pkg/r2package"
)

type Site struct {
	path string
}

func New(path string) (Site, error) {
	s := Site{path}

	// create the filesystem structure
	paths := []string{
		path,
		s.databaseSubDir(),
		s.gitSubDir(),
		s.installedSubDir(),
	}

	for _, p := range paths {
		if err := os.MkdirAll(p, 0755); err != nil {
			return Site{}, fmt.Errorf("could not create %s: %w", p, err)
		}
	}

	return s, nil
}

func (s Site) Database() database.Database {
	return database.New(s.databaseSubDir())
}

func (s Site) InstallPackage(name string) error {
	ifile, err := s.Database().GetInfoFile(name)
	if err != nil {
		return fmt.Errorf("could not find the info file: %w", err)
	}

	return s.installFromInfoFile(ifile)
}

func (s Site) InstalledPackage(name string) (r2package.InfoFile, error) {
	path := filepath.Join(s.installedSubDir(), name)
	return r2package.FromFile(path)
}

func (s Site) InstallPackageFromFile(path string) error {
	ifile, err := r2package.FromFile(path)
	if err != nil {
		return fmt.Errorf("could not read %s as a package info file: %w", path, err)
	}

	return s.installFromInfoFile(ifile)
}

func (s Site) UninstallPackage(name string) error {
	ifile, err := s.InstalledPackage(name)
	if err != nil {
		log.Print(err)
		return fmt.Errorf("could not find %s as an installed package", name)
	}

	dir, err := s.getPackageSubDir(ifile.Type)
	if err != nil {
		return fmt.Errorf("could not determine where %s is installed: %w", name, err)
	}

	installedDir := filepath.Join(dir, ifile.Name)

	if err := ifile.Uninstall(installedDir); err != nil {
		return fmt.Errorf("could not uninstall %s: %w", name, err)
	}

	return os.Remove(ifile.Path)
}

func (s Site) ListInstalledPackages() ([]r2package.Info, error) {
	dir := s.installedSubDir()

	ifiles, err := r2package.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %w", dir, err)
	}

	packages := make([]r2package.Info, 0, len(ifiles))

	for _, p := range ifiles {
		packages = append(packages, p.Info)
	}

	return packages, nil
}

func (s Site) Remove() error {
	return os.RemoveAll(s.path)
}

func (s Site) Upgrade(name string) error {
	// is the package installed?
	_, err := s.InstalledPackage(name)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("%q is not installed", name)
	}

	// is the package available in the database?
	_, err = s.Database().GetInfoFile(name)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("%q is not available in the database", name)
	}

	if err := s.UninstallPackage(name); err != nil {
		return fmt.Errorf("could not uninstall %s: %w", name, err)
	}

	if err := s.InstallPackage(name); err != nil {
		return fmt.Errorf("could not install %s: %w", name, err)
	}

	return nil
}

//
// Private
//

func (s Site) databaseSubDir() string {
	return filepath.Join(s.path, "r2pm-db")
}

func (s Site) getPackageSubDir(pkgType string) (string, error) {
	switch pkgType {
	case "git":
		return s.gitSubDir(), nil
	default:
		return "", fmt.Errorf("%q: unhandled package type", pkgType)
	}
}

func (s Site) gitSubDir() string {
	return filepath.Join(s.path, "git")
}

func (s Site) installFromInfoFile(ifile r2package.InfoFile) error {
	dir, err := s.getPackageSubDir(ifile.Type)
	if err != nil {
		return fmt.Errorf("could not determine where to install %s: %w", ifile.Name, err)
	}

	dir = filepath.Join(dir, ifile.Name)

	if err := os.Mkdir(dir, 0755); err != nil {
		return fmt.Errorf("could not create %s: %w", dir, err)
	}

	if err := ifile.Install(dir); err != nil {
		// delete the directory that we just created
		os.RemoveAll(dir)

		return fmt.Errorf("could not install %s in %s: %w", ifile.Name, dir, err)
	}

	installedFilename := filepath.Join(s.installedSubDir(), ifile.Name)

	return pkg.CopyFile(ifile.Path, installedFilename)
}

func (s Site) installedSubDir() string {
	return filepath.Join(s.path, "installed")
}
