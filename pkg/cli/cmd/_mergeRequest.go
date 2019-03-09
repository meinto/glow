package cmd

import (
	"log"

	"github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/spf13/cobra"
)

// TODO implement this cmd
func init() {
	// rootCmd.AddCommand(mergeRequestCmd)
	util.AddFlagsForMergeRequests(mergeRequestCmd)
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

		util.CheckRequiredStringField(source, "source branch")
		util.CheckRequiredStringField(target, "target branch")
		util.CheckRequiredStringField(util.MergeRequestFlags.GitlabEndpoint, "gitlab endpoint")
		util.CheckRequiredStringField(util.MergeRequestFlags.ProjectNamespace, "project namespace")
		util.CheckRequiredStringField(util.MergeRequestFlags.ProjectName, "project name")
		util.CheckRequiredStringField(util.MergeRequestFlags.GitlabCIToken, "gitlab ci token")

		util.CreateMergeRequest(source, target)
	},
}
