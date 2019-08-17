package git

import (
	"path/filepath"

	"golang.org/x/xerrors"

	"github.com/radareorg/r2pm/pkg/process"
)

const gitBin = "git"

type Repository string

func Init(path string, force bool) (Repository, error) {
	if _, err := Open(path); err == nil && !force {
		return "", xerrors.Errorf("cannot init: %s is already a git repository", path)
	}

	if err := Run([]string{"init"}, path); err != nil {
		return "", xerrors.Errorf("error while running git init in %s: %w", path, err)
	}

	return Repository(path), nil
}

func Open(path string) (Repository, error) {
	return Repository(path), Run([]string{"rev-parse"}, path)
}

func (r Repository) AddRemote(name, url string) error {
	return r.Run("remote", "add", name, url)
}

func (r Repository) Pull(remote, branch string) error {
	// Do not send the remote and branch names if they are empty
	args := []string{"pull"}

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

func Run(args []string, gitDir string) error {
	gitDir = filepath.Join(gitDir, ".git")

	_, err := process.Run(
		gitBin,
		append([]string{"--git-dir", gitDir}, args...),
		"")

	return err
}
