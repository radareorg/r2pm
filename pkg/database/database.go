package database

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"

	"github.com/radareorg/r2pm/pkg/git"
	"github.com/radareorg/r2pm/pkg/r2package"
)

const repoName = "r2pm-db"

func Delete(r2pmDir string) error {
	return os.RemoveAll(r2pmDir)
}

func Init(r2pmDir string) error {
	if err := os.MkdirAll(r2pmDir, 0755); err != nil {
		return xerrors.Errorf("could not create %s: %w", r2pmDir, err)
	}

	const repoUrl = "https://github.com/radareorg/" + repoName

	repoDir := filepath.Join(r2pmDir, repoName)

	if repo, err := git.Open(repoDir); err != nil {
		log.Printf("Cloning %s in %s", repoUrl, r2pmDir)

		args := []string{"--depth=3", "--recursive"}

		if err := git.Clone(repoUrl, r2pmDir, "", args); err != nil {
			return xerrors.Errorf("could not clone %s: %w", repoName, err)
		}
	} else {
		log.Printf("pulling the latest revision from %s", repoUrl)

		if err := repo.Run("reset", "--hard", "HEAD"); err != nil {
			return xerrors.Errorf("could not reset the repo: %w", repoName, err)
		}

		// assume origin / master
		if err := repo.Pull("", ""); err != nil {
			return xerrors.Errorf("could pull the latest revision: %w", err)
		}
	}

	validPackages := make([]string, 0)

	err := filepath.Walk(r2pmDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Read file content and parse it
		pi, err := r2package.FromFile(path)
		if err != nil {
			return nil
		}

		// Validate package
		if pi.Name != filepath.Base(path) {
			log.Printf("Invalid package name in %q: %q", path, pi.Name)
			return nil
		}

		validPackages = append(validPackages, pi.Name)

		return nil
	})

	if err != nil {
		return xerrors.Errorf("could not initialize the database: %w", err)
	}

	dbFile := filepath.Join(r2pmDir, "db.json")

	fd, err := os.Create(dbFile)
	if err != nil {
		return xerrors.Errorf("could not open %s for writing: %w", dbFile, err)
	}
	defer fd.Close()

	return json.NewEncoder(fd).Encode(validPackages)
}

func FindPackage(r2pmDir, packageName string) (*r2package.Info, error) {
	const dbSubdir = "db"

	path := filepath.Join(r2pmDir, repoName, dbSubdir, packageName)

	return r2package.FromFile(path)
}
