package cmd

import (
	"fmt"

	"github.com/meinto/glow/cmd/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(publishCmd)
	util.AddFlagsForMergeRequests(publishCmd)
}

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish a release branch",
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0] // should be semver

		source := fmt.Sprintf("release/v%s", version)
		target := "master"
		util.CreateMergeRequest(source, target)
	},
}
