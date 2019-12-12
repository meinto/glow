package cmd

import (
	"github.com/meinto/glow"

	"github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(featureCmd)
}

var featureCmd = &cobra.Command{
	Use:   "feature",
	Short: "create a feature branch",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		featureName := args[0]

		feature, err := glow.NewFeature(viper.GetString("author"), featureName)
		util.CheckForError(err, "NewFeature")

		g, err := util.GetGitClient()
		util.CheckForError(err, "GetGitClient")

		err = g.Create(feature, rootCmdOptions.SkipChecks)
		util.CheckForError(err, "Create")

		g.Checkout(feature)
		util.CheckForError(err, "Checkout")
	},
}
