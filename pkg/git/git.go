package git

import (
	"os/exec"

	"golang.org/x/xerrors"
)

type Repository string

func Init(path string, force bool) (Repository, error) {
	_, err := Open(path)
	if err == nil && !force {
		return "", xerrors.Errorf("cannot init: %s is already a git repository", path)
	}

	if err := Run([]string{"init"}, path); err != nil {
		return "", xerrors.Errorf("error while running git init in %s: %w", path, err)
	}

	return Repository(path), nil
}

func Open(path string) (Repository, error) {
	if Run([]string{"rev-parse", "--is-inside-work-tree"}, path) != nil {
		return "", xerrors.Errorf("%s is not a git repository")
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

func Clone(repoUrl, wd, dstDir string, opts []string) error {
	args := []string{"clone"}
	args = append(args, opts...)
	args = append(args, repoUrl)

	if dstDir != "" {
		args = append(args, dstDir)
	}

	return Run(args, wd)
}

func Run(args []string, wd string) error {
	cmd := exec.Command("git", args...)

	if wd != "" {
		cmd.Dir = wd
	}

	return cmd.Run()
}
