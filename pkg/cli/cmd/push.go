package cmd

import (
	"github.com/meinto/glow"
	. "github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/spf13/cobra"
)

var pushCmdOptions struct {
	AddAll        bool
	CommitMessage string
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().BoolVar(&pushCmdOptions.AddAll, "addAll", false, "add all changes made on the current branch")
	pushCmd.Flags().StringVar(&pushCmdOptions.CommitMessage, "commitMessage", "", "add a commit message (flag --addAll required)")
	AddFlagsForMergeRequests(pushCmd)
}

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push changes",
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

		if pushCmdOptions.AddAll {
			ExitOnError(g.AddAll())
			ExitOnError(g.Stash())
			ExitOnError(g.Checkout(currentBranch))
			ExitOnError(g.StashPop())
			ExitOnError(g.AddAll())

			if pushCmdOptions.CommitMessage != "" {
				ExitOnError(g.Commit(pushCmdOptions.CommitMessage))
			}
		}

		g.Push(false)
		ExitOnError(err)
	},
}
