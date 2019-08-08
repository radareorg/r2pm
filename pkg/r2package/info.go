package r2package

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/xerrors"
	"gopkg.in/yaml.v2"

	"github.com/radareorg/r2pm/pkg/git"
)

const gitSubdir = "git"

type Info struct {
	Name          string
	Type          string
	Repo          string
	Desc          string
	InstallCmds   []string `yaml:"install"`
	UninstallCmds []string `yaml:"uninstall"`
}

func FromFile(path string) (*Info, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	d := yaml.NewDecoder(fd)
	d.SetStrict(true)

	info := &Info{}

	err = d.Decode(&info)

	return info, err
}

func (i Info) Install(r2pmDir string) error {
	if i.Type != "git" {
		return xerrors.Errorf("package type %s not implemented", i.Type)
	}

	gitDir := filepath.Join(r2pmDir, gitSubdir)
	repoPath := i.InstallDir(r2pmDir)

	if _, err := git.Open(repoPath); err == nil {
		return xerrors.Errorf(
			"%s is already installed - %s exists",
			i.Name,
			repoPath)
	}

	if err := os.MkdirAll(gitDir, 0755); err != nil {
		return xerrors.Errorf("could not create %s: %w", gitDir, err)
	}

	if err := git.Clone(i.Repo, gitDir, "", nil); err != nil {
		return xerrors.Errorf("could not clone %s: %w", i.Repo, err)
	}

	for idx, line := range i.InstallCmds {
		fields := strings.Fields(line)

		cmd := exec.Command(fields[0], fields[1:]...)
		cmd.Dir = repoPath

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

func (i Info) InstallDir(r2pmDir string) string {
	gitDir := filepath.Join(r2pmDir, gitSubdir)
	return filepath.Join(gitDir, filepath.Base(i.Repo))
}

func (i Info) Uninstall(r2pmDir string) error {
	for idx, line := range i.UninstallCmds {
		fields := strings.Fields(line)

		cmd := exec.Command(fields[0], fields[1:]...)
		cmd.Dir = i.InstallDir(r2pmDir)

		if err := cmd.Run(); err != nil {
			return xerrors.Errorf(
				"uninstall command #%d [%q] failed: %w",
				idx+1,
				line,
				err)
		}

	}

	return os.RemoveAll(i.InstallDir(r2pmDir))
}
