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

func (c *Command) Do(str string, args ...interface{}) (stdout, stderr bytes.Buffer, err error) {
	formattedStr := str
	if strings.Contains(str, "%") {
		formattedStr = fmt.Sprintf(str, args...)
	}
	cmd := exec.Command(c.shell, "-c", formattedStr)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Start()
	if err != nil {
		return bytes.Buffer{}, bytes.Buffer{}, err
	}
	return stdout, stderr, cmd.Wait()
}
