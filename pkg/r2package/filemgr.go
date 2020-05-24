package r2package

import "github.com/radareorg/r2pm/pkg"

type fileMgr struct{}

func (fm *fileMgr) CopyFile(src, dst string) error {
	return pkg.CopyFile(src, dst)
}
