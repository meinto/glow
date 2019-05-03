package cmd

import (
	"github.com/meinto/glow"
	"github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/spf13/cobra"
)

var hotfixCmdOptions struct {
	Push              bool
	PostHotfixScript  string
	PostHotfixCommand []string
	VersionFile       string
	VersionFileType   string
}

func init() {
	rootCmd.AddCommand(hotfixCmd)

	hotfixCmd.Flags().BoolVar(&hotfixCmdOptions.Push, "push", false, "push created hotfix branch")
	hotfixCmd.Flags().StringVar(&hotfixCmdOptions.PostHotfixScript, "postHotfix", "", "script that executes after switching to hotfix branch")
	hotfixCmd.Flags().StringArrayVar(&hotfixCmdOptions.PostHotfixCommand, "postHotfixCommand", []string{}, "commands which should be executed after switching to hotfix branch")

	hotfixCmd.Flags().StringVar(&hotfixCmdOptions.VersionFile, "versionFile", "VERSION", "name of git-semver version file")
	hotfixCmd.Flags().StringVar(&hotfixCmdOptions.VersionFileType, "versionFileType", "raw", "git-semver version file type")
}

var hotfixCmd = &cobra.Command{
	Use:   "hotfix",
	Short: "create a hotfix branch",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		g, err := util.GetGitClient()
		util.CheckForError(err, "GetGitClient")

		version, s := util.ProcessVersion(
			args[0],
			hotfixCmdOptions.VersionFile,
			hotfixCmdOptions.VersionFileType,
		)

		hotfix, err := glow.NewHotfix(version)
		util.CheckForError(err, "NewHotfix")

		err = g.Create(hotfix, rootCmdOptions.SkipChecks)
		util.CheckForError(err, "Create")

		g.Checkout(hotfix)
		util.CheckForError(err, "Checkout")

		if util.IsSemanticVersion(args[0]) {
			err = s.SetNextVersion(args[0])
			util.CheckForError(err, "semver SetNextVersion")
		} else {
			err = s.SetVersion(version)
			util.CheckForError(err, "semver SetVersion")
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		util.PostRunWithCurrentVersion(
			hotfixCmdOptions.VersionFile,
			hotfixCmdOptions.VersionFileType,
			hotfixCmdOptions.PostHotfixScript,
			hotfixCmdOptions.PostHotfixCommand,
			hotfixCmdOptions.Push,
		)
	},
}
