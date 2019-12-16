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
		util.ExitOnError(err)

		g, err := util.GetGitClient()
		util.ExitOnError(g.Create(feature, rootCmdOptions.SkipChecks))
		util.ExitOnError(g.Checkout(feature))
	},
}
