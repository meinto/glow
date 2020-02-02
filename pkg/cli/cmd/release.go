package cmd

import (
	"github.com/meinto/glow"
	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
	"github.com/meinto/glow/pkg/cli/cmd/internal/util"
	"github.com/spf13/cobra"
)

var releaseCmdOptions struct {
	Push               bool
	PostReleaseScript  string
	PostReleaseCommand []string
	VersionFile        string
	VersionFileType    string
}

type ReleaseCommand struct {
	command.Service
}

func (cmd *ReleaseCommand) PostSetup(parent command.Service) command.Service {
	parent.Add(cmd)
	cmd.Cmd().Flags().BoolVar(&releaseCmdOptions.Push, "push", false, "push created release branch")
	cmd.Cmd().Flags().StringVar(&releaseCmdOptions.PostReleaseScript, "postRelease", "", "script that executes after switching to release branch")
	cmd.Cmd().Flags().StringArrayVar(&releaseCmdOptions.PostReleaseCommand, "postReleaseCommand", []string{}, "commands which should be executed after switching to release branch")

	cmd.Cmd().Flags().StringVar(&releaseCmdOptions.VersionFile, "versionFile", "VERSION", "name of git-semver version file")
	cmd.Cmd().Flags().StringVar(&releaseCmdOptions.VersionFileType, "versionFileType", "raw", "git-semver version file type")
	return cmd
}

var ReleaseCmd = SetupReleaseCommand(RootCmd)

func SetupReleaseCommand(parent command.Service) command.Service {
	return command.Setup(&ReleaseCommand{
		&command.Command{
			Command: &cobra.Command{
				Use:   "release",
				Short: "create a release branch",
				Args:  cobra.MinimumNArgs(1),
			},
			Run: func(cmd command.Service, args []string) {
				version := util.ProcessVersionS(args[0], cmd.SemverClient())

				release, err := glow.NewRelease(version)
				util.ExitOnError(err)

				util.ExitOnError(cmd.GitClient().Create(release, rootCmdOptions.SkipChecks))

				_, _, err = cmd.GitClient().Checkout(release)
				util.ExitOnError(err)

				if util.IsSemanticVersion(args[0]) {
					util.ExitOnError(cmd.SemverClient().SetNextVersion(args[0]))
				} else {
					util.ExitOnError(cmd.SemverClient().SetVersion(version))
				}
			},
			PostRun: func(cmd command.Service, args []string) {
				util.PostRunWithCurrentVersionS(
					cmd.SemverClient(),
					cmd.GitClient(),
					releaseCmdOptions.PostReleaseScript,
					releaseCmdOptions.PostReleaseCommand,
					releaseCmdOptions.Push,
				)
			},
		},
	}).PostSetup(parent)
}
