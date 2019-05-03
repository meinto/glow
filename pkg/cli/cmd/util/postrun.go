package util

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

func PostRunWithVersion(
	versionArg, versionFile, versionFileType, postReleaseScript string,
	postReleaseCommand []string,
	push bool,
) {
	g, err := GetGitClient()
	CheckForError(err, "GetGitClient")

	version, _ := ProcessVersion(versionArg, versionFile, versionFileType)

	if postReleaseScript != "" {
		postRelease(version, postReleaseScript)
	}
	if len(postReleaseCommand) > 0 {
		for _, command := range postReleaseCommand {
			execute(version, command)
		}
	}

	if push {
		err = g.AddAll()
		CheckForError(err, "AddAll")

		err = g.Commit("[glow] Add post release changes")
		CheckForError(err, "Commit")

		err = g.Push(true)
		CheckForError(err, "Push")
	}
}

func postRelease(version, script string) {
	pathToFile, err := filepath.Abs(script)
	if err != nil {
		log.Println("cannot find post-release script", err)
	}
	cmd := exec.Command(pathToFile, version)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Println("error while executing post-release script", err)
	}
	log.Println("post release:")
	log.Println(out.String())
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
	if err != nil {
		log.Println("error while executing post-release script", err.Error(), stderr.String())
	}
	log.Println("post release:")
	log.Println(stdout.String())
}