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
		util.ExitOnError(err)

		// err := g.Fetch()
		// util.CheckForError(err, "Fetch")

		gp, err := util.GetGitProvider()
		util.ExitOnError(err)

		var currentBranch glow.Branch
		if rootCmdOptions.CI {
			cb, err := gp.GetCIBranch()
			util.ExitOnError(err)
			currentBranch = cb
		} else {
			cb, _, _, err := g.CurrentBranch()
			util.ExitOnError(err)
			currentBranch = cb
		}

		err = gp.Close(currentBranch)
		if !util.MergeRequestFlags.Gracefully {
			util.ExitOnError(err)
		} else {
			log.Println(err)
		}
	},
}
