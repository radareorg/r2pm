package r2package

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type executor interface {
	Run(c *exec.Cmd) error
}

type fetcher interface {
	Fetch(ctx context.Context, dir string) error
}

type fileManager interface {
	CopyFile(src, dst string) error
}

type Installer struct {
	cmdExecutor  executor
	dirs         R2Dirs
	fileManager  fileManager
	getFetcher   func(Source) (fetcher, error)
	logger       *log.Logger
	tmpDirGetter func(string, string) (string, error)
}

type R2Dirs struct {
	// Headers is the directory where r2 headers are located
	Headers string

	// Libs is the directory where r2 libraries are located
	Libs string

	// Plugins is the destination directory for plugins
	Plugins string
}

func NewInstaller(m Manifest, logger *log.Logger, dirs R2Dirs) *Installer {
	return &Installer{
		cmdExecutor:  &commandExecutor{},
		dirs:         dirs,
		getFetcher:   getFetcher,
		logger:       logger,
		tmpDirGetter: ioutil.TempDir,
	}
}

func (i *Installer) Install(ctx context.Context, m Manifest) error {
	const goos = runtime.GOOS

	osFamily, err := osFromGOOS(goos)
	if err != nil {
		return fmt.Errorf("could not determine the OS family for GOOS=%q: %v", goos, err)
	}

	instructions, ok := m.Install[osFamily]
	if !ok {
		return fmt.Errorf("no install instructions for %q (GOOS: %q)", osFamily, goos)
	}

	f, err := i.getFetcher(instructions.Source)
	if err != nil {
		return fmt.Errorf("invalid source: %v", err)
	}

	tmpDir, err := i.tmpDirGetter("", "r2pm-*")
	if err != nil {
		return fmt.Errorf("could not create a temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	i.logger.Printf("Fetching in %s", tmpDir)

	if err := f.Fetch(ctx, tmpDir); err != nil {
		return fmt.Errorf("could not fetch sources in %s: %v", tmpDir, err)
	}

	for _, cmd := range instructions.Commands {
		items := strings.Split(cmd, " ")

		cmd := exec.CommandContext(ctx, items[0], items[1:]...)
		cmd.Dir = tmpDir

		if err := i.cmdExecutor.Run(cmd); err != nil {
			return fmt.Errorf("error while executing %q: %v", cmd, err)
		}
	}

	for _, outFile := range instructions.Out {
		src := filepath.Join(tmpDir, outFile.Path)
		dst := filepath.Join(i.dirs.Plugins, outFile.Path)

		if err := i.fileManager.CopyFile(src, dst); err != nil {
			return fmt.Errorf("could not copy %s to %s: %v", src, dst, err)
		}
	}

	return nil
}
