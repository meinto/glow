package util

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var MergeRequestFlags struct {
	Githost          string
	ProjectNamespace string
	ProjectName      string
	GitlabCIToken    string
}

func AddFlagsForMergeRequests(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&MergeRequestFlags.Githost, "endpoint", "e", "", "git host endpoint")
	cmd.Flags().StringVarP(&MergeRequestFlags.ProjectNamespace, "namespace", "n", "", "project namespace")
	cmd.Flags().StringVarP(&MergeRequestFlags.ProjectName, "project", "p", "", "project name")
	cmd.Flags().StringVarP(&MergeRequestFlags.GitlabCIToken, "token", "t", "", "gitlab ci token")
	viper.BindPFlag("githost", cmd.Flags().Lookup("endpoint"))
	viper.BindPFlag("projectNamespace", cmd.Flags().Lookup("namespace"))
	viper.BindPFlag("projectName", cmd.Flags().Lookup("project"))
	viper.BindPFlag("gitlabCIToken", cmd.Flags().Lookup("token"))
}
