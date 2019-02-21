package cmd

import (
	"log"

	"github.com/meinto/glow/cmd/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var mergeRequestCmdOptions struct {
	GitlabEndpoint   string
	ProjectNamespace string
	ProjectName      string
	GitlabCIToken    string
}

func init() {
	rootCmd.AddCommand(mergeRequestCmd)
	mergeRequestCmd.Flags().StringVarP(&mergeRequestCmdOptions.GitlabEndpoint, "endpoint", "e", "", "gitlab endpoint")
	mergeRequestCmd.Flags().StringVarP(&mergeRequestCmdOptions.ProjectNamespace, "namespace", "n", "", "project namespace")
	mergeRequestCmd.Flags().StringVarP(&mergeRequestCmdOptions.ProjectName, "project", "p", "", "project name")
	mergeRequestCmd.Flags().StringVarP(&mergeRequestCmdOptions.GitlabCIToken, "token", "t", "", "gitlab ci token")
	viper.BindPFlag("gitlabEndpoint", mergeRequestCmd.Flags().Lookup("endpoint"))
	viper.BindPFlag("projectNamespace", mergeRequestCmd.Flags().Lookup("namespace"))
	viper.BindPFlag("projectName", mergeRequestCmd.Flags().Lookup("project"))
	viper.BindPFlag("gitlabCIToken", mergeRequestCmd.Flags().Lookup("token"))
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
		utils.CheckRequiredStringField(mergeRequestCmdOptions.GitlabEndpoint, "gitlab endpoint")
		utils.CheckRequiredStringField(mergeRequestCmdOptions.ProjectNamespace, "project namespace")
		utils.CheckRequiredStringField(mergeRequestCmdOptions.ProjectName, "project name")
		utils.CheckRequiredStringField(mergeRequestCmdOptions.GitlabCIToken, "gitlab ci token")

		utils.CreateMergeRequest(source, target)
	},
}
