package r2package

import (
	"os"
	"os/exec"
	"strings"

	"golang.org/x/xerrors"

	"github.com/radareorg/r2pm/pkg/git"
)

type gitInstaller struct {
	info Info
}

func (g gitInstaller) install(inDir string) error {
	const (
		remoteName   = "origin"
		remoteBranch = "master"
	)

	repo, err := git.Init(inDir, false)
	if err != nil {
		return xerrors.Errorf("could not init the repository: %w", err)
	}

	if err := repo.AddRemote(remoteName, g.info.Repo); err != nil {
		return xerrors.Errorf("could not add the remote: %w", err)
	}

	if err := repo.Pull(remoteName, remoteBranch, []string{"--depth=1"}); err != nil {
		return xerrors.Errorf("could not git pull: %w", err)
	}

	for idx, line := range g.info.InstallCmds {
		fields := strings.Fields(line)

		cmd := exec.Command(fields[0], fields[1:]...)
		cmd.Dir = inDir

		if err := cmd.Run(); err != nil {
			return xerrors.Errorf(
				"install command #%d [%q] failed: %w",
				idx+1,
				line,
				err)
		}
	}

	return nil
}

func (g gitInstaller) uninstall(fromDir string) error {
	for idx, line := range g.info.UninstallCmds {
		fields := strings.Fields(line)

		cmd := exec.Command(fields[0], fields[1:]...)
		cmd.Dir = fromDir

		if err := cmd.Run(); err != nil {
			return xerrors.Errorf(
				"uninstall command #%d [%q] failed: %w",
				idx+1,
				line,
				err)
		}

	}

	return os.RemoveAll(fromDir)
}
