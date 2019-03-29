package cmd

import (
	"github.com/meinto/glow"
	"github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/meinto/glow/semver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
				viper.GetString("gitPath"),
				releaseCmdOptions.VersionFile,
				releaseCmdOptions.VersionFileType,
			)
			v, err := s.GetCurrentVersion(args[0])
			util.CheckForError(err, "semver GetCurrentVersion")
			version = v
		}

		gp, err := util.GetGitProvider()
		util.CheckForError(err, "GetGitProvider")

		currentBranch, err := glow.NewRelease(version)
		util.CheckForError(err, "NewRelease")

		gp.Close(currentBranch)
		util.CheckForError(err, "Close")
	},
}
