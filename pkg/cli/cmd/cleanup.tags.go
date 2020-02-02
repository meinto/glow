package cmd

import (
	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
	"github.com/meinto/glow/pkg/cli/cmd/internal/util"
	"github.com/spf13/cobra"
)

type CleanupTagsCommand struct {
	command.Service
}

func (cmd *CleanupTagsCommand) PostSetup(parent command.Service) command.Service {
	parent.Add(cmd)
	return cmd
}

var cleanupTagsCmd = SetupCleanupTagsCommand(cleanupCmd)

func SetupCleanupTagsCommand(parent command.Service) command.Service {
	return command.Setup(&CleanupTagsCommand{
		&command.Command{
			Command: &cobra.Command{
				Use:   "tags",
				Short: "cleanup tags",
			},
			Run: func(cmd command.Service, args []string) {
				util.ExitOnError(cmd.GitClient().CleanupTags(cleanupCmdFlags.cleanupUntracked))
			},
		},
	}).PostSetup(parent)
}
