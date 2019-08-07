package database

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/xerrors"

	"github.com/radareorg/r2pm/pkg/r2package"
)

func Init(r2pmDir string) error {
	if err := os.MkdirAll(r2pmDir, 0755); err != nil {
		return xerrors.Errorf("could not create %s: %w", r2pmDir, err)
	}

	const (
		repoName = "r2pm-db"
		repoUrl  = "https://github.com/radareorg/" + repoName
	)

	repoDir := filepath.Join(r2pmDir, repoName)

	// is repoDir already a git repository?
	cmdRepoExists := exec.Command(
		"git",
		"rev-parse",
		"--is-inside-work-tree")

	cmdRepoExists.Dir = repoDir

	if err := cmdRepoExists.Run(); err != nil {
		log.Printf("Cloning %s in %s", repoUrl, r2pmDir)

		args := []string{"git", "clone", "--depth=3", "--recursive"}

		if err := runGit(args, r2pmDir); err != nil {
			return xerrors.Errorf("could not clone %s: %w", repoName, err)
		}
	} else {
		log.Printf("pulling the latest revision from %s", repoUrl)

		if err := runGit([]string{"reset", "--hard", "HEAD"}, repoDir); err != nil {
			return xerrors.Errorf("could not reset the repo: %w", repoName, err)
		}

		if err := runGit([]string{"pull"}, repoDir); err != nil {
			return xerrors.Errorf("could pull the latest revision: %w", repoName, err)
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

func runGit(args []string, wd string) error {
	cmd := exec.Command("git", args...)

	if wd != "" {
		cmd.Dir = wd
	}

	return cmd.Run()
}
