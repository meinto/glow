package cmd

import (
	"github.com/meinto/glow"
	l "github.com/meinto/glow/logging"
	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
	"github.com/meinto/glow/pkg/cli/cmd/internal/util"
	"github.com/meinto/glow/semver"
	"github.com/spf13/cobra"
)

type CloseReleaseCommand struct {
	command.Service
}

func (cmd *CloseReleaseCommand) PostSetup(parent command.Service) command.Service {
	parent.Add(cmd)
	util.AddFlagsForMergeRequests(cmd.Cmd())
	return cmd
}

var closeReleaseCmd = SetupCloseReleaseCommand(RootCmd)

func SetupCloseReleaseCommand(parent command.Service) command.Service {
	return command.Setup(&CloseReleaseCommand{
		&command.Command{
			Command: &cobra.Command{
				Use:   "release",
				Short: "close a release branch",
				Args:  cobra.MinimumNArgs(1),
			},
			Run: func(cmd command.Service, args []string) {
				version := args[0]

				if version == "current" {
					pathToRepo, _, _, err := cmd.GitClient().GitRepoPath()
					util.ExitOnError(err)

					s := semver.NewSemverService(
						pathToRepo,
						"/bin/bash",
						releaseCmdOptions.VersionFile,
						releaseCmdOptions.VersionFileType,
					)
					v, err := s.GetCurrentVersion()
					util.ExitOnError(err)
					version = v
				}

				currentBranch, err := glow.NewRelease(version)
				util.ExitOnError(err)

				err = cmd.GitProvider().Close(currentBranch)
				if !util.MergeRequestFlags.Gracefully {
					util.ExitOnError(err)
				} else {
					l.Log().Error(err)
				}
			},
		},
	})
}
