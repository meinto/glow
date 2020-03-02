package cmd

import (
	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
	"github.com/meinto/glow/pkg/cli/cmd/internal/util"
	"github.com/spf13/cobra"
)

type CleanupBranchesCommand struct {
	command.Service
}

func (cmd *CleanupBranchesCommand) PostSetup(parent command.Service) command.Service {
	parent.Add(cmd)
	return cmd
}

var cleanupBranchesCmd = SetupCleanupBranchesCommand(cleanupCmd)

func SetupCleanupBranchesCommand(parent command.Service) command.Service {
	return command.Setup(&CleanupBranchesCommand{
		&command.Command{
			Command: &cobra.Command{
				Use:   "branches",
				Short: "cleanup branches",
			},
			Run: func(cmd command.Service, args []string) {
				util.ExitOnError(cmd.GitClient().CleanupBranches(cleanupCmdFlags.cleanupGone, cleanupCmdFlags.cleanupUntracked))
			},
		},
	}, parent)
}
