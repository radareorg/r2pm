package r2package

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Info struct {
	Name      string   `yaml:"name"`
	Type      string   `yaml:"type"`
	Repo      string   `yaml:"repo"`
	Desc      string   `yaml:"desc"`
	Install   []string `yaml:"install"`
	Uninstall []string `yaml:"uninstall"`
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
