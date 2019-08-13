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
	res, err := process.Run(gitBin, []string{"rev-parse", "--show-toplevel"}, path)
	if err != nil {
		return "", xerrors.Errorf("%s is not a git repository", path)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", xerrors.Errorf("could not determine the absolute path for %s: %w", path, err)
	}

	// TODO does this work on Windows? git probably prints CRLF there.
	if res.Stdout.String() != (absPath + "\n") {
		return "", xerrors.Errorf("%s is not a git repository", path)
	}

	return Repository(path), nil
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

func Run(args []string, wd string) error {
	_, err := process.Run(gitBin, args, wd)
	return err
}
