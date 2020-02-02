package cmd

import (
	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
	"github.com/spf13/cobra"
)

var cleanupCmdFlags struct {
	cleanupGone      bool
	cleanupUntracked bool
}

type CleanupCommand struct {
	command.Service
}

func (cmd *CleanupCommand) PostSetup(parent command.Service) command.Service {
	parent.Add(cmd)
	cmd.Cmd().PersistentFlags().BoolVar(&cleanupCmdFlags.cleanupGone, "gone", false, "cleanup branches which are gone on remote")
	cmd.Cmd().PersistentFlags().BoolVar(&cleanupCmdFlags.cleanupUntracked, "untracked", false, "cleanup branches which are gone on remote")
	return cmd
}

var cleanupCmd = SetupCleanupCommand(RootCmd)

func SetupCleanupCommand(parent command.Service) command.Service {
	return command.Setup(&CleanupCommand{
		&command.Command{
			Command: &cobra.Command{
				Use:   "cleanup",
				Short: "cleanup branches",
			},
			Run: func(cmd command.Service, args []string) {},
		},
	}).PostSetup(parent)
}
