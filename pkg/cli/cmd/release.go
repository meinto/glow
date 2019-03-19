package cmd

import (
	"bytes"
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
}

func init() {
	rootCmd.AddCommand(releaseCmd)
	releaseCmd.Flags().BoolVar(&releaseCmdOptions.Push, "push", false, "push created release branch")
	releaseCmd.Flags().StringVar(&releaseCmdOptions.PostReleaseScript, "postRelease", "", "script that executes after switching to release branch")
	releaseCmd.Flags().StringArrayVar(&releaseCmdOptions.PostReleaseCommand, "postReleaseCommand", []string{}, "commands which should be executed after switching to release branch")
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
			s = semver.NewGitSemverService(pathToRepo, viper.GetString("gitPath"))
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

		if releaseCmdOptions.Push {
			g.Push(true)
			util.CheckForError(err, "Push")
		}

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
	cmdParts := strings.Split(command, " ")
	args := []string{"-c"}
	args = append(args, cmdParts...)
	cmd := exec.Command("/bin/bash", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println("error while executing post-release script", err)
	}
	log.Println("post release:")
	log.Println(out.String())
}
