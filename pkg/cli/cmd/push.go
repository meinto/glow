package cmd

import (
	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
	"github.com/meinto/glow/pkg/cli/cmd/internal/util"
	"github.com/spf13/cobra"
)

var pushCmdOptions struct {
	AddAll        bool
	CommitMessage string
}

type PushCommand struct {
	command.Service
}

func (cmd *PushCommand) PostSetup(parent command.Service) command.Service {
	parent.Add(cmd)
	cmd.Cmd().Flags().BoolVar(&pushCmdOptions.AddAll, "addAll", false, "add all changes made on the current branch")
	cmd.Cmd().Flags().StringVar(&pushCmdOptions.CommitMessage, "commitMessage", "", "add a commit message (flag --addAll required)")
	util.AddFlagsForMergeRequests(cmd.Cmd())
	return cmd
}

var pushCmd = SetupPushCommand(RootCmd)

func SetupPushCommand(parent command.Service) command.Service {
	return command.Setup(&PushCommand{
		&command.Command{
			Command: &cobra.Command{
				Use:   "push",
				Short: "push changes",
			},
			Run: func(cmd command.Service, args []string) {
				currentBranch := cmd.CurrentBranch(RootCmdOptions.CI)

				if pushCmdOptions.AddAll {
					util.ExitOnError(cmd.GitClient().AddAll())
					util.ExitOnError(cmd.GitClient().Stash())
					util.ExitOnError(cmd.GitClient().Checkout(currentBranch))
					util.ExitOnError(cmd.GitClient().StashPop())
					util.ExitOnError(cmd.GitClient().AddAll())

					if pushCmdOptions.CommitMessage != "" {
						util.ExitOnError(cmd.GitClient().Commit(pushCmdOptions.CommitMessage))
					}
				}

				exists, _, _, err := cmd.GitClient().RemoteBranchExists(currentBranch.BranchName())
				util.ExitOnError(err)

				_, _, err = cmd.GitClient().Push(!exists)
				util.ExitOnError(err)
			},
		},
	}, parent)
}
