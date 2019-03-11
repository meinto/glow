package cmd

import (
	"bytes"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/meinto/glow"
	"github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/spf13/cobra"
)

var releaseCmdOptions struct {
	PostReleaseScript string
}

func init() {
	rootCmd.AddCommand(releaseCmd)
	releaseCmd.Flags().StringVar(&releaseCmdOptions.PostReleaseScript, "postRelease", "", "script that executes after switching to release branch")
}

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "create a release branch",
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]

		release, err := glow.NewRelease(version)
		util.CheckForError(err, "NewRelease")

		g, err := util.GetGitClient()
		util.CheckForError(err, "GetGitClient")

		err = g.Create(release)
		util.CheckForError(err, "Create")

		g.Checkout(release)
		util.CheckForError(err, "Checkout")

		if releaseCmdOptions.PostReleaseScript != "" {
			postRelease(version)
		}
	},
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
