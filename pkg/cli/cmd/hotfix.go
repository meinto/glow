package cmd

import (
	"github.com/meinto/glow"
	"github.com/meinto/glow/git"
	"github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(hotfixCmd)
}

var hotfixCmd = &cobra.Command{
	Use:   "hotfix",
	Short: "create a hotfix branch",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hotixName := args[0]

		hotfix, err := glow.NewHotfix(viper.GetString("author"), hotixName)
		util.CheckForError(err, "NewHotfix")

		g := git.NewGoGitService()

		err = g.Create(hotfix)
		util.CheckForError(err, "Create")

		g.Checkout(hotfix)
		util.CheckForError(err, "Checkout")
	},
}
