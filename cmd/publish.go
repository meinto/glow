package cmd

import (
	"fmt"

	"github.com/meinto/glow/cmd/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(publishCmd)
	utils.AddFlagsForMergeRequests(publishCmd)
}

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish a release branch",
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0] // should be semver

		source := fmt.Sprintf("release/v%s", version)
		target := "master"
		utils.CreateMergeRequest(source, target)
	},
}
