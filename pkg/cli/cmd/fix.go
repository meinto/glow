package cmd

import (
	"github.com/meinto/glow"
	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
	"github.com/meinto/glow/pkg/cli/cmd/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type FixCommand struct {
	command.Service
}

func (cmd *FixCommand) PostSetup(parent command.Service) command.Service {
	parent.Add(cmd)
	return cmd
}

var fixCmd = SetupFixCommand(RootCmd)

func SetupFixCommand(parent command.Service) command.Service {
	return command.Setup(&FixCommand{
		&command.Command{
			Command: &cobra.Command{
				Use:   "fix",
				Short: "create a fix branch",
				Args:  cobra.MinimumNArgs(1),
			},
			Run: func(cmd command.Service, args []string) {
				fixName := args[0]

				fix, err := glow.NewFix(viper.GetString("author"), fixName)
				util.ExitOnError(err)

				util.ExitOnError(cmd.GitClient().Create(fix, rootCmdOptions.SkipChecks))
				util.ExitOnError(cmd.GitClient().Checkout(fix))
			},
		},
	}, parent)
}
