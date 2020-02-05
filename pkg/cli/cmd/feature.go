package cmd

import (
	"github.com/meinto/glow"

	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
	"github.com/meinto/glow/pkg/cli/cmd/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type FeatureCommand struct {
	command.Service
}

func (cmd *FeatureCommand) PostSetup(parent command.Service) command.Service {
	parent.Add(cmd)
	return cmd
}

var featureCmd = SetupFeatureCommand(RootCmd)

func SetupFeatureCommand(parent command.Service) command.Service {
	return command.Setup(&FeatureCommand{
		&command.Command{
			Command: &cobra.Command{
				Use:   "feature",
				Short: "create a feature branch",
				Args:  cobra.MinimumNArgs(1),
			},
			Run: func(cmd command.Service, args []string) {
				featureName := args[0]

				feature, err := glow.NewFeature(viper.GetString("author"), featureName)
				util.ExitOnError(err)

				util.ExitOnError(cmd.GitClient().Create(feature, rootCmdOptions.SkipChecks))
				util.ExitOnError(cmd.GitClient().Checkout(feature))
			},
		},
	}, parent)
}
