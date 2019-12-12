package cmd

import (
	"github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pushCmd)
	util.AddFlagsForMergeRequests(pushCmd)
}

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push changes",
	Run: func(cmd *cobra.Command, args []string) {

		g, err := util.GetGitClient()
		util.CheckForError(err, "GetGitClient")

		err = g.Push(false)
		util.CheckForError(err, "Push")
	},
}
