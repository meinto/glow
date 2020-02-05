package cmd

import (
	"github.com/meinto/glow"
	l "github.com/meinto/glow/logging"
	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
	"github.com/meinto/glow/pkg/cli/cmd/internal/util"
	"github.com/spf13/cobra"
)

type CloseCommand struct {
	command.Service
}

func (cmd *CloseCommand) PostSetup(parent command.Service) command.Service {
	parent.Add(cmd)
	util.AddFlagsForMergeRequests(cmd.Cmd())
	return cmd
}

var closeCmd = SetupCloseCommand(RootCmd)

func SetupCloseCommand(parent command.Service) command.Service {
	return command.Setup(&CloseCommand{
		&command.Command{
			Command: &cobra.Command{
				Use:   "close",
				Short: "close a branch",
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

				err := cmd.GitProvider().Close(currentBranch)
				if !util.MergeRequestFlags.Gracefully {
					util.ExitOnError(err)
				} else {
					l.Log().Error(err)
				}
			},
		},
	}, parent)
}
