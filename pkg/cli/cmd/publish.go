package cmd

import (
	"github.com/meinto/glow"
	l "github.com/meinto/glow/logging"
	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
	"github.com/meinto/glow/pkg/cli/cmd/internal/util"
	"github.com/spf13/cobra"
)

type PublishCommand struct {
	command.Service
}

func (cmd *PublishCommand) PostSetup(parent command.Service) command.Service {
	parent.Add(cmd)
	util.AddFlagsForMergeRequests(cmd.Cmd())
	return cmd
}

var publishCmd = SetupPublishCommand(RootCmd)

func SetupPublishCommand(parent command.Service) command.Service {
	return command.Setup(&PublishCommand{
		&command.Command{
			Command: &cobra.Command{
				Use:   "publish",
				Short: "publish a release branch",
			},
			Run: func(cmd command.Service, args []string) {
				var currentBranch glow.Branch
				if rootCmdOptions.CI {
					cb, err := cmd.GitProvider().GetCIBranch()
					util.ExitOnError(err)
					currentBranch = cb
				} else {
					cb, _, _, err := cmd.GitClient().CurrentBranch()
					util.ExitOnError(err)
					currentBranch = cb
				}

				err := cmd.GitProvider().Publish(currentBranch)
				if !util.MergeRequestFlags.Gracefully {
					util.ExitOnError(err)
				} else {
					l.Log().Error(err)
				}
			},
		},
	})
}
