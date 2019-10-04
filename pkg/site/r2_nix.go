// +build darwin freebsd linux openbsd

package site

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/xerrors"

	"github.com/radareorg/r2pm/pkg/git"
	"github.com/radareorg/r2pm/pkg/process"
)

func (s Site) InstallRadare2(prefix string) error {
	srcDir := filepath.Join(s.gitSubDir(), "radare2")

	if err := os.MkdirAll(srcDir, 0755); err != nil {
		return xerrors.Errorf(
			"could not create the filesystem tree for %s: %v",
			srcDir,
			err)
	}

	log.Print("Opening " + srcDir)

	repo, err := git.Open(srcDir)
	if err != nil {
		log.Printf("Could not open %s as a git repo: %v", srcDir, err)
		log.Print("Running git init")

		if repo, err = git.Init(srcDir, false); err != nil {
			return xerrors.Errorf("could not run git init: %v", err)
		}

		origin := "https://github.com/radare/radare2"

		log.Print("Setting the origin to " + origin)
		if err = repo.AddRemote("origin", origin); err != nil {
			return xerrors.Errorf("could not set origin: %v", err)
		}
	}

	if err := repo.Pull("origin", "master", []string{"--depth=1"}); err != nil {
		return err
	}

	// Allow ./configure to be executed
	configurePath := filepath.Join(srcDir, "configure")

	log.Print("Allowing the execution of " + configurePath)

	if err := os.Chmod(configurePath, 0755); err != nil {
		return err
	}

	env := os.Environ()
	makeBin := "make"

	if runtime.GOOS == "freebsd" {
		makeBin = "gmake"
		env = append(env, "CC=clang")
	}

	cmdConfigure := exec.Command("./configure", "--prefix="+prefix)
	cmdConfigure.Dir = srcDir
	cmdConfigure.Env = env

	log.Print("Running " + strings.Join(cmdConfigure.Args, " "))

	if output, err := cmdConfigure.CombinedOutput(); err != nil {
		log.Print(output)
		return err
	}

	cmdMake := exec.Command(makeBin)
	cmdMake.Dir = srcDir
	cmdMake.Env = env

	log.Print("Running " + strings.Join(cmdMake.Args, " "))

	if out, err := cmdMake.Output(); err != nil {
		log.Print(out)
		return err
	}

	if _, err := process.Run(makeBin, []string{"install"}, srcDir); err != nil {
		return err
	}

	return nil
}
