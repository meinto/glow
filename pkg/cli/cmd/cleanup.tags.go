package cmd

import (
	"github.com/meinto/glow/pkg/cli/cmd/internal/util"
	"github.com/spf13/cobra"
)

func init() {
	cleanupCmd.Cmd().AddCommand(cleanupTagsCmd)
}

var cleanupTagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "cleanup tags",
	Run: func(cmd *cobra.Command, args []string) {

		g, err := util.GetGitClient()
		util.ExitOnError(err)

		util.ExitOnError(g.CleanupTags(cleanupCmdFlags.cleanupUntracked))
	},
}
