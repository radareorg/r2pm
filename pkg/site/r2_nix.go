// +build darwin freebsd linux openbsd

package site

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/radareorg/r2pm/pkg/git"
	"github.com/radareorg/r2pm/pkg/process"
)

func (s Site) InstallRadare2(prefix, version string) error {
	srcDir := filepath.Join(s.gitSubDir(), "radare2")

	if err := os.MkdirAll(srcDir, 0755); err != nil {
		return fmt.Errorf(
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
			return fmt.Errorf("could not run git init: %v", err)
		}

		origin := "https://github.com/radareorg/radare2"

		log.Print("Setting the origin to " + origin)
		if err = repo.AddRemote("origin", origin); err != nil {
			return fmt.Errorf("could not set origin: %v", err)
		}
	}

	if err := repo.Fetch(); err != nil {
		return fmt.Errorf("Could not fetch: %v", err)
	}

	if err := repo.Checkout(version); err != nil {
		return fmt.Errorf("could not checkout %q: %v", version, err)
	}

	if err := repo.Pull("origin", version, nil); err != nil {
		return err
	}

	// Allow ./configure to be executed
	configurePath := filepath.Join(srcDir, "configure")

	log.Print("Allowing the execution of " + configurePath)

	if err := os.Chmod(configurePath, 0755); err != nil {
		return err
	}

	env := os.Environ()

	cmdConfigure := exec.Command("./configure", "--prefix="+prefix)
	cmdConfigure.Dir = srcDir
	cmdConfigure.Env = env

	log.Print("Running " + strings.Join(cmdConfigure.Args, " "))

	if out, err := cmdConfigure.CombinedOutput(); err != nil {
		log.Print(string(out))
		return err
	}

	makeBin := "make"

	if runtime.GOOS == "freebsd" {
		makeBin = "gmake"
	}

	cmdMake := exec.Command(makeBin)
	cmdMake.Dir = srcDir
	cmdMake.Env = env

	log.Print("Running " + strings.Join(cmdMake.Args, " "))

	if out, err := cmdMake.Output(); err != nil {
		log.Print(string(out))
		return err
	}

	if _, err := process.Run(makeBin, []string{"install"}, srcDir); err != nil {
		return err
	}

	return nil
}
