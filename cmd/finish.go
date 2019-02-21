package cmd

import (
	"fmt"

	"github.com/meinto/glow/cmd/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(finishCmd)
	utils.AddFlagsForMergeRequests(finishCmd)
}

var finishCmd = &cobra.Command{
	Use:   "finish",
	Short: "finish a release branch",
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0] // should be semver

		source := fmt.Sprintf("release/v%s", version)
		target := "develop"
		utils.CreateMergeRequest(source, target)
	},
}
