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
		util.ExitOnError(err)
		util.ExitOnError(g.CleanupBranches(cleanupCmdFlags.cleanupGone, cleanupCmdFlags.cleanupUntracked))
	},
}
