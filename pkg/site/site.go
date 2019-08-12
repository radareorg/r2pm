package site

import (
	"os"
	"path/filepath"

	"golang.org/x/xerrors"

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
			return Site{}, xerrors.Errorf("could not create %s: %w")
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
		return xerrors.Errorf("could not find the info file: %w", err)
	}

	dir, err := s.getPackageSubDir(ifile.Type)
	if err != nil {
		return xerrors.Errorf("could not determine where to install %s: %w", name, err)
	}

	dir = filepath.Join(dir, ifile.Name)

	if err := os.Mkdir(dir, 0755); err != nil {
		return xerrors.Errorf("could not create %s: %w", dir, err)
	}

	if err := ifile.Install(dir); err != nil {
		// delete the directory that we just created
		os.RemoveAll(dir)

		return xerrors.Errorf("could not install %s in %s: %w", name, dir, err)
	}

	installedFilename := filepath.Join(s.installedSubDir(), ifile.Name)

	return pkg.CopyFile(ifile.Path, installedFilename)
}

func (s Site) UninstallPackage(name string) error {
	installedInfoFile := filepath.Join(s.installedSubDir(), name)

	ifile, err := r2package.FromFile(installedInfoFile)
	if err != nil {
		return xerrors.Errorf("could not find %s as an installed package: %w", name, err)
	}

	dir, err := s.getPackageSubDir(ifile.Type)
	if err != nil {
		return xerrors.Errorf("could not determine where %s is installed: %w", name, err)
	}

	installedDir := filepath.Join(dir, ifile.Name)

	if err := ifile.Uninstall(installedDir); err != nil {
		return xerrors.Errorf("could not uninstall %s: %w", name, err)
	}

	return os.Remove(installedInfoFile)
}

func (s Site) ListInstalledPackages() ([]r2package.Info, error) {
	dir := s.installedSubDir()

	ifiles, err := r2package.ReadDir(dir)
	if err != nil {
		return nil, xerrors.Errorf("could not read %s: %w", dir, err)
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
		return "", xerrors.Errorf("%q: unhandled package type", pkgType)
	}
}

func (s Site) gitSubDir() string {
	return filepath.Join(s.path, "git")
}

func (s Site) installedSubDir() string {
	return filepath.Join(s.path, "installed")
}
