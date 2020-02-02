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

func SetupReleaseCmd(parent, cmd command.Service) command.Service {
	parent.Add(cmd)
	cmd.Cmd().Flags().BoolVar(&releaseCmdOptions.Push, "push", false, "push created release branch")
	cmd.Cmd().Flags().StringVar(&releaseCmdOptions.PostReleaseScript, "postRelease", "", "script that executes after switching to release branch")
	cmd.Cmd().Flags().StringArrayVar(&releaseCmdOptions.PostReleaseCommand, "postReleaseCommand", []string{}, "commands which should be executed after switching to release branch")

	cmd.Cmd().Flags().StringVar(&releaseCmdOptions.VersionFile, "versionFile", "VERSION", "name of git-semver version file")
	cmd.Cmd().Flags().StringVar(&releaseCmdOptions.VersionFileType, "versionFileType", "raw", "git-semver version file type")
	return cmd
}

func ReleaseCmd(parent command.Service) command.Service {
	return SetupReleaseCmd(parent, &command.Command{
		Command: &cobra.Command{
			Use:   "release",
			Short: "create a release branch",
			Args:  cobra.MinimumNArgs(1),
		},
		Run: func(cmd command.Service, args []string) {
			pathToRepo, _, _, err := cmd.GitClient().GitRepoPath()
			util.ExitOnError(err)

			version, s := util.ProcessVersion(
				args[0],
				releaseCmdOptions.VersionFile,
				releaseCmdOptions.VersionFileType,
				pathToRepo,
			)

			release, err := glow.NewRelease(version)
			util.ExitOnError(err)

			util.ExitOnError(cmd.GitClient().Create(release, rootCmdOptions.SkipChecks))

			_, _, err = cmd.GitClient().Checkout(release)
			util.ExitOnError(err)

			if util.IsSemanticVersion(args[0]) {
				util.ExitOnError(s.SetNextVersion(args[0]))
			} else {
				util.ExitOnError(s.SetVersion(version))
			}
		},
		PostRun: func(cmd command.Service, args []string) {
			util.PostRunWithCurrentVersion(
				releaseCmdOptions.VersionFile,
				releaseCmdOptions.VersionFileType,
				releaseCmdOptions.PostReleaseScript,
				releaseCmdOptions.PostReleaseCommand,
				releaseCmdOptions.Push,
			)
		},
	})
}
