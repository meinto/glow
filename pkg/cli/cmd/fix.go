package cmd

import (
	"github.com/meinto/glow"
	. "github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(fixCmd)
}

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "create a fix branch",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fixName := args[0]

		fix, err := glow.NewFix(viper.GetString("author"), fixName)
		ExitOnError(err)

		g, err := GetGitClient()
		ExitOnError(err)

		ExitOnError(g.Create(fix, rootCmdOptions.SkipChecks))
		ExitOnError(g.Checkout(fix))
	},
}
