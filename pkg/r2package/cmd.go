package r2package

import "os/exec"

type commandExecutor struct{}

func (ce *commandExecutor) Run(cmd *exec.Cmd) error {
	return cmd.Run()
}
