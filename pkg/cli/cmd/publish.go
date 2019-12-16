package cmd

import (
	"log"

	"github.com/meinto/glow"
	"github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(publishCmd)
	util.AddFlagsForMergeRequests(publishCmd)
}

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish a release branch",
	Run: func(cmd *cobra.Command, args []string) {

		g, err := util.GetGitClient()
		util.ExitOnError(err)

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

		err = gp.Publish(currentBranch)
		if !util.MergeRequestFlags.Gracefully {
			util.ExitOnError(err)
		} else {
			log.Println(err)
		}
	},
}
