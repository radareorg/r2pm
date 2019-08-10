package git

import (
	"os/exec"

	"golang.org/x/xerrors"
)

type Repository string

func Open(path string) (Repository, error) {
	if Run([]string{"rev-parse", "--is-inside-work-tree"}, path) != nil {
		return "", xerrors.Errorf("%s is not a git repository")
	}

	return Repository(path), nil
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
