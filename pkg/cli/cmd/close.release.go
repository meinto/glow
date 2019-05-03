package cmd

import (
	"log"

	"github.com/meinto/glow"
	"github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/meinto/glow/semver"
	"github.com/spf13/cobra"
)

func init() {
	closeCmd.AddCommand(closeReleaseCmd)
	util.AddFlagsForMergeRequests(closeReleaseCmd)
}

var closeReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "close a release branch",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]

		if version == "current" {
			g, err := util.GetGitClient()
			util.CheckForError(err, "GetGitClient")

			pathToRepo, err := g.GitRepoPath()
			util.CheckForError(err, "semver GitRepoPath")

			s := semver.NewSemverService(
				pathToRepo,
				"/bin/bash",
				releaseCmdOptions.VersionFile,
				releaseCmdOptions.VersionFileType,
			)
			v, err := s.GetCurrentVersion()
			util.CheckForError(err, "semver GetCurrentVersion")
			version = v
		}

		gp, err := util.GetGitProvider()
		util.CheckForError(err, "GetGitProvider")

		currentBranch, err := glow.NewRelease(version)
		util.CheckForError(err, "NewRelease")

		err = gp.Close(currentBranch)
		if !util.MergeRequestFlags.Gracefully {
			util.CheckForError(err, "Close")
		} else {
			log.Println(err)
		}
	},
}
