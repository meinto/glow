package cmd

import (
	"log"

	"github.com/meinto/glow"
	"github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(closeCmd)
	util.AddFlagsForMergeRequests(closeCmd)
}

var closeCmd = &cobra.Command{
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

		err = gp.Close(currentBranch)
		if !util.MergeRequestFlags.Gracefully {
			util.CheckForError(err, "Close")
		} else {
			log.Println(err)
		}
	},
}
