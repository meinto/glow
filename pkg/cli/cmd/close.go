package cmd

import (
	"github.com/meinto/glow"
	"github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(finishCmd)
	util.AddFlagsForMergeRequests(finishCmd)
}

var finishCmd = &cobra.Command{
	Use:   "close",
	Short: "close a branch",
	Run: func(cmd *cobra.Command, args []string) {

		g, err := util.GetGitClient()
		util.CheckForError(err, "GetGitClient")

		// err := g.Fetch()
		// util.CheckForError(err, "Fetch")

		gp, err := util.GetGitProvider()
		util.CheckForError(err, "GetGitProvider")

		var currentBranch glow.Branch
		if rootCmdOptions.CI {
			cb, err := gp.GetCIBranch()
			util.CheckForError(err, "CurrentBranch")
			currentBranch = cb
		} else {
			cb, err := g.CurrentBranch()
			util.CheckForError(err, "CurrentBranch")
			currentBranch = cb
		}

		gp.Close(currentBranch)
		util.CheckForError(err, "Close")
	},
}
