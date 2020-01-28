package util

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var MergeRequestFlags struct {
	GitProviderDomain string
	GitProvider       string
	ProjectNamespace  string
	ProjectName       string
	Token             string
	Gracefully        bool
}

func AddFlagsForMergeRequests(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&MergeRequestFlags.GitProviderDomain, "endpoint", "e", "", "git host endpoint")
	cmd.Flags().StringVar(&MergeRequestFlags.GitProvider, "gitProvider", "gitlab", "git provider e.g.: gitlab, github, bitbucket")
	cmd.Flags().StringVarP(&MergeRequestFlags.ProjectNamespace, "namespace", "n", "", "project namespace")
	cmd.Flags().StringVarP(&MergeRequestFlags.ProjectName, "project", "p", "", "project name")
	cmd.Flags().StringVarP(&MergeRequestFlags.Token, "token", "t", "", "gitlab ci token")
	cmd.Flags().BoolVar(&MergeRequestFlags.Gracefully, "gracefully", false, "only log on error")
	viper.BindPFlag("gitProviderDomain", cmd.Flags().Lookup("endpoint"))
	viper.BindPFlag("gitProvider", cmd.Flags().Lookup("gitProvider"))
	viper.BindPFlag("projectNamespace", cmd.Flags().Lookup("namespace"))
	viper.BindPFlag("projectName", cmd.Flags().Lookup("project"))
	viper.BindPFlag("token", cmd.Flags().Lookup("token"))
}
