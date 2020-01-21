package cmd

import (
	"github.com/meinto/glow"
	. "github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/spf13/cobra"
)

var releaseCmdOptions struct {
	Push               bool
	PostReleaseScript  string
	PostReleaseCommand []string
	VersionFile        string
	VersionFileType    string
}

func init() {
	rootCmd.AddCommand(releaseCmd)
	releaseCmd.Flags().BoolVar(&releaseCmdOptions.Push, "push", false, "push created release branch")
	releaseCmd.Flags().StringVar(&releaseCmdOptions.PostReleaseScript, "postRelease", "", "script that executes after switching to release branch")
	releaseCmd.Flags().StringArrayVar(&releaseCmdOptions.PostReleaseCommand, "postReleaseCommand", []string{}, "commands which should be executed after switching to release branch")

	releaseCmd.Flags().StringVar(&releaseCmdOptions.VersionFile, "versionFile", "VERSION", "name of git-semver version file")
	releaseCmd.Flags().StringVar(&releaseCmdOptions.VersionFileType, "versionFileType", "raw", "git-semver version file type")
}

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "create a release branch",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		g, err := GetGitClient()
		ExitOnError(err)

		version, s := ProcessVersion(
			args[0],
			releaseCmdOptions.VersionFile,
			releaseCmdOptions.VersionFileType,
		)

		release, err := glow.NewRelease(version)
		ExitOnError(err)

		ExitOnError(g.Create(release, rootCmdOptions.SkipChecks))

		_, _, err = g.Checkout(release)
		ExitOnError(err)

		if IsSemanticVersion(args[0]) {
			ExitOnError(s.SetNextVersion(args[0]))
		} else {
			ExitOnError(s.SetVersion(version))
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		PostRunWithCurrentVersion(
			releaseCmdOptions.VersionFile,
			releaseCmdOptions.VersionFileType,
			releaseCmdOptions.PostReleaseScript,
			releaseCmdOptions.PostReleaseCommand,
			releaseCmdOptions.Push,
		)
	},
}
