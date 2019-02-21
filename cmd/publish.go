package cmd

import (
	"fmt"

	"github.com/meinto/glow/cmd/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var publishCmdOptions struct {
	GitlabEndpoint   string
	ProjectNamespace string
	ProjectName      string
	GitlabCIToken    string
}

func init() {
	rootCmd.AddCommand(publishCmd)
	publishCmd.Flags().StringVarP(&publishCmdOptions.GitlabEndpoint, "endpoint", "e", "", "gitlab endpoint")
	publishCmd.Flags().StringVarP(&publishCmdOptions.ProjectNamespace, "namespace", "n", "", "project namespace")
	publishCmd.Flags().StringVarP(&publishCmdOptions.ProjectName, "project", "p", "", "project name")
	publishCmd.Flags().StringVarP(&publishCmdOptions.GitlabCIToken, "token", "t", "", "gitlab ci token")
	viper.BindPFlag("gitlabEndpoint", publishCmd.Flags().Lookup("endpoint"))
	viper.BindPFlag("projectNamespace", publishCmd.Flags().Lookup("namespace"))
	viper.BindPFlag("projectName", publishCmd.Flags().Lookup("project"))
	viper.BindPFlag("gitlabCIToken", publishCmd.Flags().Lookup("token"))
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
