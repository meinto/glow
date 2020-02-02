package cmd

import (
	"github.com/meinto/glow"
	"github.com/meinto/glow/pkg/cli/cmd/internal/util"
	"github.com/spf13/cobra"
)

var pushCmdOptions struct {
	AddAll        bool
	CommitMessage string
}

func init() {
	RootCmd.Cmd().AddCommand(pushCmd)
	pushCmd.Flags().BoolVar(&pushCmdOptions.AddAll, "addAll", false, "add all changes made on the current branch")
	pushCmd.Flags().StringVar(&pushCmdOptions.CommitMessage, "commitMessage", "", "add a commit message (flag --addAll required)")
	util.AddFlagsForMergeRequests(pushCmd)
}

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push changes",
	Run: func(cmd *cobra.Command, args []string) {

		g, err := util.GetGitClient()
		util.ExitOnError(err)

		gp, err := util.GetGitProvider()
		util.ExitOnError(err)

		var currentBranch glow.Branch
		if rootCmdOptions.CI {
			cb := gp.GetCIBranch()
			util.ExitOnError(err)
			currentBranch = cb
		} else {
			cb, _, _, err := g.CurrentBranch()
			util.ExitOnError(err)
			currentBranch = cb
		}

		if pushCmdOptions.AddAll {
			util.ExitOnError(g.AddAll())
			util.ExitOnError(g.Stash())
			util.ExitOnError(g.Checkout(currentBranch))
			util.ExitOnError(g.StashPop())
			util.ExitOnError(g.AddAll())

			if pushCmdOptions.CommitMessage != "" {
				util.ExitOnError(g.Commit(pushCmdOptions.CommitMessage))
			}
		}

		g.Push(false)
		util.ExitOnError(err)
	},
}
