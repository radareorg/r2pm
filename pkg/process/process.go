package process

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

type Result struct {
	Stderr bytes.Buffer
	Stdout bytes.Buffer
}

func Run(binary string, args []string, wd string) (Result, error) {
	res := Result{}

	cmd := exec.Command(binary, args...)
	cmd.Stderr = &res.Stderr
	cmd.Stdout = &res.Stdout

	if wd != "" {
		cmd.Dir = wd
	} else {
		wd = "."
	}

	log.Printf("Running %q in %s", strings.Join(cmd.Args, " "), wd)

	return res, cmd.Run()
}
