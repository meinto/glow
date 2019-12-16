package cmd

import (
	"bytes"
	"os/exec"
)

type Cmd struct {
	cmd *exec.Cmd
}

func (c *Cmd) Run() (stdout, stderr string, err error) {
	var stdoutBuff, stderrBuff bytes.Buffer
	c.cmd.Stdout = &stdoutBuff
	c.cmd.Stderr = &stderrBuff
	err = c.cmd.Run()
	return stdoutBuff.String(), stderrBuff.String(), err
}

func (c *Cmd) Get() *exec.Cmd {
	return c.cmd
}
