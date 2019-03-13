package cmd

import (
	"github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/spf13/cobra"
)

var cleanupCmdFlags struct {
	cleanupGone      bool
	cleanupUntracked bool
}

func init() {
	rootCmd.AddCommand(cleanupCmd)
	cleanupCmd.PersistentFlags().BoolVar(&cleanupCmdFlags.cleanupGone, "gone", false, "cleanup branches which are gone on remote")
	cleanupCmd.PersistentFlags().BoolVar(&cleanupCmdFlags.cleanupUntracked, "untracked", false, "cleanup branches which are gone on remote")
}

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "cleanup branches",
	Run: func(cmd *cobra.Command, args []string) {

		g, err := util.GetGitClient()
		util.CheckForError(err, "GetGitClient")

		err = g.CleanupBranches(cleanupCmdFlags.cleanupGone, cleanupCmdFlags.cleanupUntracked)
		util.CheckForError(err, "Close")
	},
}
