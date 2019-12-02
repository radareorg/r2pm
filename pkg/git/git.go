package git

import (
	"fmt"
	"path/filepath"

	"github.com/radareorg/r2pm/pkg/process"
)

const gitBin = "git"

type Repository string

func Init(path string, force bool) (Repository, error) {
	if _, err := Open(path); err == nil && !force {
		return "", fmt.Errorf("cannot init: %s is already a git repository", path)
	}

	if err := Run([]string{"init"}, path); err != nil {
		return "", fmt.Errorf("error while running git init in %s: %w", path, err)
	}

	return Repository(path), nil
}

func Open(path string) (Repository, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf(
			"could not get the absolute path for %s: %w",
			path,
			err)
	}

	args := []string{"--git-dir", filepath.Join(absPath, ".git"), "rev-parse"}

	// Check that absPath contains a .git Repository
	return Repository(path), Run(args, "")
}

func (r Repository) AddRemote(name, url string) error {
	return r.Run("remote", "add", name, url)
}

func (r Repository) Checkout(ref string) error {
	return r.Run("checkout", ref)
}

func (r Repository) Fetch() error {
	return r.Run("fetch")
}

func (r Repository) Pull(remote, branch string, opts []string) error {
	// Do not send the remote and branch names if they are empty
	args := append([]string{"pull"}, opts...)

	if remote != "" {
		args = append(args, remote)
	}

	if branch != "" {
		args = append(args, branch)
	}

	return r.Run(args...)
}

func (r Repository) Run(args ...string) error {
	return Run(args, string(r))
}

func Run(args []string, wd string) error {
	_, err := process.Run(gitBin, args, wd)
	return err
}
