package database

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"

	"github.com/radareorg/r2pm/pkg/git"
	"github.com/radareorg/r2pm/pkg/r2package"
)

const (
	dbSubdir = "db"
	repoName = "r2pm-db"
)

type Database struct {
	path string
}

func New(path string) Database {
	return Database{path}
}

func (d Database) InitOrUpdate() error {
	const (
		remoteName   = "origin"
		remoteBranch = "master"
	)

	repo, err := git.Open(d.path)
	if err != nil {
		// Create the repo if it does not exist
		repo, err = git.Init(d.path, false)
		if err != nil {
			return xerrors.Errorf("could not initialize the database repo: %w", err)
		}

		if err := repo.AddRemote(remoteName, "https://github.com/radareorg/"+repoName); err != nil {
			return xerrors.Errorf("could not add the remote: %w", err)
		}
	}

	// assume origin / master
	if err := repo.Pull(remoteName, remoteBranch); err != nil {
		return xerrors.Errorf("could not pull the latest revision: %w", err)
	}

	return nil
}

func (d Database) Delete() error {
	return os.RemoveAll(d.path)
}

func (d Database) GetInfoFile(packageName string) (r2package.InfoFile, error) {
	path := filepath.Join(d.path, dbSubdir, packageName)

	return r2package.FromFile(path)
}

// ListAvailablePackages returns a slice of strings containing the names of all the installer packages.
func (d Database) ListAvailablePackages() ([]string, error) {
	dirs, err := ioutil.ReadDir(filepath.Join(d.path, dbSubdir))
	if err != nil {
		return nil, xerrors.Errorf("could not list the directory: %w", err)
	}

	packages := make([]string, 0, len(dirs))

	for _, dir := range dirs {
		// skip all except directories
		if !dir.IsDir() {
			continue
		}

		packages = append(packages, filepath.Base(dir.Name()))
	}

	return packages, nil
}
