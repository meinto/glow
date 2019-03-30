package testenv

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Command struct {
	shell string
}

func NewCommand() *Command {
	return &Command{"/bin/bash"}
}

func (c *Command) Do(str string, args ...interface{}) (bytes.Buffer, error) {
	formattedStr := str
	if strings.Contains(str, "%") {
		formattedStr = fmt.Sprintf(str, args...)
	}
	cmd := exec.Command(c.shell, "-c", formattedStr)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	return stdout, cmd.Run()
}
