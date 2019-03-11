package cmd

import (
	"github.com/meinto/glow/githost"
	"github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(publishCmd)
	util.AddFlagsForMergeRequests(publishCmd)
}

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish a release branch",
	Run: func(cmd *cobra.Command, args []string) {

		g, err := util.GetGitClient()
		util.CheckForError(err, "GetGitClient")

		// err := g.Fetch()
		// util.CheckForError(err, "Fetch")

		currentBranch, err := g.CurrentBranch()
		util.CheckForError(err, "CurrentBranch")

		gh := githost.NewGitlabService(
			viper.GetString("gitlabEndpoint"),
			viper.GetString("projectNamespace"),
			viper.GetString("projectName"),
			viper.GetString("gitlabCIToken"),
		)
		gh = githost.NewLoggingService(logger, gh)

		gh.Publish(currentBranch)
		util.CheckForError(err, "Close")
	},
}
