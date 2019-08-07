package r2package

import "golang.org/x/xerrors"

type Info struct {
	Name      string   `yaml:"name"`
	Type      string   `yaml:"type"`
	Repo      string   `yaml:"repo"`
	Desc      string   `yaml:"desc"`
	Install   []string `yaml:"install"`
	Uninstall []string `yaml:"uninstall"`
}

func FromFile(path string) (*Info, error) {
	return nil, xerrors.New("not implemented")
}
