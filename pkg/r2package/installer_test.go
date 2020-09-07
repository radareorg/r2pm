//go:generate go run github.com/golang/mock/mockgen -package r2package -source installer.go -destination installer_mock_test.go executor,fetcher,fileManager

package r2package

import (
	"context"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestInstaller_Install(t *testing.T) {
	ctrl := gomock.NewController(t)

	f := NewMockfetcher(ctrl)
	e := NewMockexecutor(ctrl)
	fm := NewMockfileManager(ctrl)

	const (
		dirPlugins = "/a/plugin/dir"
		dirTmp     = "/a/tmp/dir"

		pathOutFile1 = "path/to/file/1"
		pathOutFile2 = "path/to/file/2"
	)

	getFetcher := func(_ Source) (fetcher, error) {
		return f, nil
	}

	tmpDirGetter := func(_, _ string) (string, error) {
		return dirTmp, nil
	}

	i := &Installer{
		cmdExecutor:  e,
		dirs:         R2Dirs{Plugins: dirPlugins},
		fileManager:  fm,
		getFetcher:   getFetcher,
		logger:       log.New(ioutil.Discard, "", 0), // TODO: figure out how to mock this
		tmpDirGetter: tmpDirGetter,
	}

	instructions := InstallInstructions{
		Commands: []string{"command 1", "command 2"},
		Out: []ManagedFile{
			{
				Path: pathOutFile1,
				Type: "exe",
			},
			{
				Path: pathOutFile2,
				Type: "shared-lib",
			},
		},
	}

	m := Manifest{
		Name:        "test-package",
		Version:     "1.2.3",
		Description: "A test package",
		Install: map[string]InstallInstructions{
			"freebsd": instructions,
			"linux":   instructions,
			"macos":   instructions,
			"windows": instructions,
		},
	}

	// Make sure we always pass the same context
	ctx := context.WithValue(context.TODO(), "test-time", time.Now().UnixNano())

	cmd1 := exec.CommandContext(ctx, "command", "1")
	cmd1.Dir = dirTmp

	cmd2 := exec.CommandContext(ctx, "command", "2")
	cmd2.Dir = dirTmp

	srcFile1 := filepath.Join(dirTmp, pathOutFile1)
	dstFile1 := filepath.Join(dirPlugins, pathOutFile1)

	srcFile2 := filepath.Join(dirTmp, pathOutFile2)
	dstFile2 := filepath.Join(dirPlugins, pathOutFile2)

	gomock.InOrder(
		f.EXPECT().Fetch(ctx, dirTmp),
		e.EXPECT().Run(cmd1),
		e.EXPECT().Run(cmd2),
		fm.EXPECT().CopyFile(srcFile1, dstFile1),
		fm.EXPECT().CopyFile(srcFile2, dstFile2),
	)

	if err := i.Install(ctx, m); err != nil {
		t.Fatal(err)
	}
}
