package cmd

import (
	"github.com/meinto/glow"
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
		g, err := util.GetGitClient()
		util.ExitOnError(err)

		version, s := util.ProcessVersion(
			args[0],
			releaseCmdOptions.VersionFile,
			releaseCmdOptions.VersionFileType,
		)

		release, err := glow.NewRelease(version)
		util.ExitOnError(err)

		util.ExitOnError(g.Create(release, rootCmdOptions.SkipChecks))

		_, _, err = g.Checkout(release)
		util.ExitOnError(err)

		if util.IsSemanticVersion(args[0]) {
			util.ExitOnError(s.SetNextVersion(args[0]))
		} else {
			util.ExitOnError(s.SetVersion(version))
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		util.PostRunWithCurrentVersion(
			releaseCmdOptions.VersionFile,
			releaseCmdOptions.VersionFileType,
			releaseCmdOptions.PostReleaseScript,
			releaseCmdOptions.PostReleaseCommand,
			releaseCmdOptions.Push,
		)
	},
}
