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

func (c *Command) Do(str string, args ...interface{}) error {
	formattedStr := str
	if strings.Contains(str, "%") {
		formattedStr = fmt.Sprintf(str, args...)
	}
	cmd := exec.Command(c.shell, "-c", formattedStr)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	return cmd.Run()
}
