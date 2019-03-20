package cmd

import (
	"github.com/meinto/glow"
	"github.com/meinto/glow/pkg/cli/cmd/util"
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

		gp, err := util.GetGitProvider()
		util.CheckForError(err, "GetGitProvider")

		currentBranch, err := glow.NewRelease(version)
		util.CheckForError(err, "NewRelease")

		gp.Close(currentBranch)
		util.CheckForError(err, "Close")
	},
}
