package cmd

import (
	"github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/spf13/cobra"
)

func init() {
	cleanupCmd.AddCommand(cleanupTagsCmd)
}

var cleanupTagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "cleanup tags",
	Run: func(cmd *cobra.Command, args []string) {

		g, err := util.GetGitClient()
		util.CheckForError(err, "GetGitClient")

		err = g.CleanupTags(cleanupCmdFlags.cleanupUntracked)
		util.CheckForError(err, "CleanupTags")
	},
}
