package util

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var MergeRequestFlags struct {
	GitlabEndpoint   string
	ProjectNamespace string
	ProjectName      string
	GitlabCIToken    string
}

func AddFlagsForMergeRequests(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&MergeRequestFlags.GitlabEndpoint, "endpoint", "e", "", "gitlab endpoint")
	cmd.Flags().StringVarP(&MergeRequestFlags.ProjectNamespace, "namespace", "n", "", "project namespace")
	cmd.Flags().StringVarP(&MergeRequestFlags.ProjectName, "project", "p", "", "project name")
	cmd.Flags().StringVarP(&MergeRequestFlags.GitlabCIToken, "token", "t", "", "gitlab ci token")
	viper.BindPFlag("gitlabEndpoint", cmd.Flags().Lookup("endpoint"))
	viper.BindPFlag("projectNamespace", cmd.Flags().Lookup("namespace"))
	viper.BindPFlag("projectName", cmd.Flags().Lookup("project"))
	viper.BindPFlag("gitlabCIToken", cmd.Flags().Lookup("token"))
}
