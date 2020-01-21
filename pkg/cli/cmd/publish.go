package cmd

import (
	"log"

	"github.com/meinto/glow"
	. "github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(publishCmd)
	AddFlagsForMergeRequests(publishCmd)
}

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish a release branch",
	Run: func(cmd *cobra.Command, args []string) {

		g, err := GetGitClient()
		ExitOnError(err)

		gp, err := GetGitProvider()
		ExitOnError(err)

		var currentBranch glow.Branch
		if rootCmdOptions.CI {
			cb := gp.GetCIBranch()
			ExitOnError(err)
			currentBranch = cb
		} else {
			cb, _, _, err := g.CurrentBranch()
			ExitOnError(err)
			currentBranch = cb
		}

		err = gp.Publish(currentBranch)
		if !MergeRequestFlags.Gracefully {
			ExitOnError(err)
		} else {
			log.Println(err)
		}
	},
}
