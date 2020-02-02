package util

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/meinto/glow/git"
	l "github.com/meinto/glow/logging"
	"github.com/meinto/glow/semver"
)

func PostRunWithCurrentVersion(
	versionFile, versionFileType, postReleaseScript string,
	postReleaseCommand []string,
	push bool,
) {
	g, err := GetGitClient()
	ExitOnError(err)

	pathToRepo, _, _, err := g.GitRepoPath()
	ExitOnError(err)

	version, _ := ProcessVersion("current", versionFile, versionFileType, pathToRepo)

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

func PostRunWithCurrentVersionS(
	semverClient semver.Service,
	gitClient git.Service,
	postReleaseScript string,
	postReleaseCommand []string,
	push bool,
) {

	if postReleaseScript != "" {
		version := ProcessVersionS("current", semverClient)
		postRelease(version, postReleaseScript)
	}
	if len(postReleaseCommand) > 0 {
		version := ProcessVersionS("current", semverClient)
		for _, command := range postReleaseCommand {
			execute(version, command)
		}
	}

	if push {
		ExitOnError(gitClient.AddAll())
		ExitOnError(gitClient.Commit("[glow] Add post release changes"))
		ExitOnError(gitClient.Push(true))
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
