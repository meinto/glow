package cmd

import (
	"github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/spf13/cobra"
)

func init() {
	cleanupCmd.AddCommand(cleanupBranchesCmd)
}

var cleanupBranchesCmd = &cobra.Command{
	Use:   "branches",
	Short: "cleanup branches",
	Run: func(cmd *cobra.Command, args []string) {

		g, err := util.GetGitClient()
		util.CheckForError(err, "GetGitClient")

		err = g.CleanupBranches(cleanupCmdFlags.cleanupGone, cleanupCmdFlags.cleanupUntracked)
		util.CheckForError(err, "CleanupBranches")
	},
}
