package cmd

import (
	"github.com/meinto/glow"

	"github.com/meinto/glow/git"
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

		g := git.NewGoGitService()

		err = g.Create(feature)
		util.CheckForError(err, "Create")

		g.Checkout(feature)
		util.CheckForError(err, "Checkout")
	},
}
