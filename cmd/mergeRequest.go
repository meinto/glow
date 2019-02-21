package cmd

import (
	"log"

	"github.com/meinto/glow/cmd/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(mergeRequestCmd)
	utils.AddFlagsForMergeRequests(mergeRequestCmd)
}

var mergeRequestCmd = &cobra.Command{
	Use:   "mergeRequest",
	Short: "create a merge request on gitlab",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Fatal("Please provide source and target branch")
		}

		source := args[0]
		target := args[1]

		utils.CheckRequiredStringField(source, "source branch")
		utils.CheckRequiredStringField(target, "target branch")
		utils.CheckRequiredStringField(utils.MergeRequestFlags.GitlabEndpoint, "gitlab endpoint")
		utils.CheckRequiredStringField(utils.MergeRequestFlags.ProjectNamespace, "project namespace")
		utils.CheckRequiredStringField(utils.MergeRequestFlags.ProjectName, "project name")
		utils.CheckRequiredStringField(utils.MergeRequestFlags.GitlabCIToken, "gitlab ci token")

		utils.CreateMergeRequest(source, target)
	},
}
