package cmd

import (
	"github.com/meinto/glow/githost"

	"github.com/meinto/glow/git"
	"github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(finishCmd)
	util.AddFlagsForMergeRequests(finishCmd)
}

var finishCmd = &cobra.Command{
	Use:   "close",
	Short: "close a branch",
	Run: func(cmd *cobra.Command, args []string) {

		g := git.NewGoGitService()
		g = git.NewLoggingService(logger, g)

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

		gh.Close(currentBranch)
		util.CheckForError(err, "Close")
	},
}
