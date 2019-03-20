package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/meinto/glow"
	"github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/meinto/glow/semver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var releaseCmdOptions struct {
	Push               bool
	PostReleaseScript  string
	PostReleaseCommand []string
	VersionFile        string
	VersionFileType    string
}

func init() {
	rootCmd.AddCommand(releaseCmd)
	releaseCmd.Flags().BoolVar(&releaseCmdOptions.Push, "push", false, "push created release branch")
	releaseCmd.Flags().StringVar(&releaseCmdOptions.PostReleaseScript, "postRelease", "", "script that executes after switching to release branch")
	releaseCmd.Flags().StringArrayVar(&releaseCmdOptions.PostReleaseCommand, "postReleaseCommand", []string{}, "commands which should be executed after switching to release branch")

	releaseCmd.Flags().StringVar(&releaseCmdOptions.VersionFile, "versionFile", "VERSION", "name of git-semver version file")
	releaseCmd.Flags().StringVar(&releaseCmdOptions.VersionFileType, "versionFileType", "raw", "git-semver version file type")
}

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "create a release branch",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]

		g, err := util.GetGitClient()
		util.CheckForError(err, "GetGitClient")

		var s semver.Service
		if hasSemverConfig() && isSemanticVersion(args[0]) {
			pathToRepo, err := g.GitRepoPath()
			util.CheckForError(err, "semver GitRepoPath")
			s = semver.NewSemverService(
				pathToRepo,
				viper.GetString("gitPath"),
				releaseCmdOptions.VersionFile,
				releaseCmdOptions.VersionFileType,
			)
			v, err := s.GetNextVersion(args[0])
			util.CheckForError(err, "semver GetNextVersion")
			version = v
		}

		release, err := glow.NewRelease(version)
		util.CheckForError(err, "NewRelease")

		err = g.Create(release)
		util.CheckForError(err, "Create")

		g.Checkout(release)
		util.CheckForError(err, "Checkout")

		if hasSemverConfig() && isSemanticVersion(args[0]) {
			err = s.SetNextVersion(args[0])
			util.CheckForError(err, "semver SetNextVersion")
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		version := args[0]
		if releaseCmdOptions.PostReleaseScript != "" {
			postRelease(version)
		}
		if len(releaseCmdOptions.PostReleaseCommand) > 0 {
			for _, command := range releaseCmdOptions.PostReleaseCommand {
				execute(version, command)
			}
		}

		if releaseCmdOptions.Push {
			g, err := util.GetGitClient()
			util.CheckForError(err, "GetGitClient")

			err = g.Push(true)
			util.CheckForError(err, "Push")
		}
	},
}

func hasSemverConfig() bool {
	if _, err := os.Stat(viper.GetString("gitPath") + "/semver.config.json"); os.IsNotExist(err) {
		return false
	}
	return true
}

func isSemanticVersion(version string) bool {
	if version == "major" || version == "minor" || version == "patch" {
		return true
	}
	return false
}

func postRelease(version string) {
	pathToFile, err := filepath.Abs(releaseCmdOptions.PostReleaseScript)
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
