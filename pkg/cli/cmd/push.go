package cmd

import (
	"github.com/meinto/glow"
	"github.com/meinto/glow/pkg/cli/cmd/util"
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
	util.AddFlagsForMergeRequests(pushCmd)
}

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push changes",
	Run: func(cmd *cobra.Command, args []string) {

		g, err := util.GetGitClient()
		util.CheckForError(err, "GetGitClient")

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

		if pushCmdOptions.AddAll {
			err = g.AddAll()
			util.CheckForError(err, "AddAll")

			err = g.Stash()
			util.CheckForError(err, "Stash")

			err = g.Checkout(currentBranch)
			util.CheckForError(err, "Checkout")

			_, _, err := g.StashPop()
			util.CheckForError(err, "StashPop")

			err = g.AddAll()
			util.CheckForError(err, "AddAll")

			if pushCmdOptions.CommitMessage != "" {
				err = g.Commit(pushCmdOptions.CommitMessage)
				util.CheckForError(err, "Commit")
			}
		}

		err = g.Push(false)
		util.CheckForError(err, "Push")
	},
}
