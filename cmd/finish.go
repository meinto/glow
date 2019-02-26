package cmd

import (
	"fmt"

	"github.com/meinto/glow/cmd/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(finishCmd)
	util.AddFlagsForMergeRequests(finishCmd)
}

var finishCmd = &cobra.Command{
	Use:   "finish",
	Short: "finish a release branch",
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0] // should be semver

		source := fmt.Sprintf("release/v%s", version)
		target := "develop"
		util.CreateMergeRequest(source, target)
	},
}
