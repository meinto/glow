package util

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	l "github.com/meinto/glow/logging"
)

func PostRunWithCurrentVersion(
	versionFile, versionFileType, postReleaseScript string,
	postReleaseCommand []string,
	push bool,
) {
	g, err := GetGitClient()
	ExitOnError(err)

	version, _ := ProcessVersion("current", versionFile, versionFileType)

	if postReleaseScript != "" {
		postRelease(version, postReleaseScript)
	}
	if len(postReleaseCommand) > 0 {
		for _, command := range postReleaseCommand {
			execute(version, command)
		}
	}

	if push {
		ExitOnError(g.AddAll())
		ExitOnError(g.Commit("[glow] Add post release changes"))
		ExitOnError(g.Push(true))
	}
}

func postRelease(version, script string) {
	pathToFile, err := filepath.Abs(script)
	l.Log().Error(err)
	cmd := exec.Command(pathToFile, version)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	l.Log().
		Stdout(out.String()).
		Error(err)
}

func execute(version, command string) {
	cmdString := string(command)
	if strings.Contains(command, "%s") {
		cmdString = fmt.Sprintf(command, version)
	}
	cmd := exec.Command("/bin/bash", "-c", cmdString)
	var stdout, stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	err := cmd.Run()
	l.Log().
		Stdout(stdout.String()).
		Stderr(stderr.String(), err)
}
