package r2package

import (
	"fmt"
	"testing"
)

func TestFromYAML(t *testing.T) {

}

func TestPackage_Check(t *testing.T) {
	t.Parallel()

	valid := Package{
		Name:        "some-package",
		Description: "some description",
		Version:     "1.2.3",
		Install: map[string]InstallInstructions{
			"linux": {},
		},
	}

	t.Run("valid example", func(t *testing.T) {
		if err := valid.Check(); err != nil {
			t.Fatal(err)
		}
	})

	for _, n := range []string{"", "some package", "some/package"} {
		t.Run(fmt.Sprintf("Name: %q", n), func(t *testing.T) {
			p := valid
			p.Name = n

			if err := p.Check(); err == nil {
				t.Fail()
			}
		})
	}

	t.Run("Empty Description", func(t *testing.T) {
		p := valid
		p.Description = ""

		if err := p.Check(); err == nil {
			t.Fail()
		}
	})

	for _, v := range []string{"", "1", "1.", "1.2", "1.2.", "1.2.3-1"} {
		t.Run(fmt.Sprintf("Version: %q", v), func(t *testing.T) {
			p := valid
			p.Version = v

			if err := p.Check(); err == nil {
				t.Fail()
			}
		})
	}

	t.Run("empty installation instructions", func(t *testing.T) {
		p := valid
		p.Install = make(map[string]InstallInstructions)

		if err := p.Check(); err == nil {
			t.Fail()
		}
	})
}
