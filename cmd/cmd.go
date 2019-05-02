package cmd

import "os/exec"

type CmdExecutor interface {
	Command(str string) *exec.Cmd
}

type cmdExecutor struct {
	executor string
	dir      string
}

func NewCmdExecutor(executor string) CmdExecutor {
	return &cmdExecutor{executor, "."}
}

func NewCmdExecutorInDir(executor, dir string) CmdExecutor {
	return &cmdExecutor{executor, dir}
}

func (c *cmdExecutor) Command(str string) *exec.Cmd {
	cmd := exec.Command(c.executor, "-c", str)
	if c.dir != "." {
		cmd.Dir = c.dir
	}
	return cmd
}
