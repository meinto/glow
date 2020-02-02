package cmd

import (
	"github.com/meinto/glow"
	"github.com/meinto/glow/pkg/cli/cmd/internal/util"
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
	RootCmd.Cmd().AddCommand(hotfixCmd)

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
		util.ExitOnError(err)

		version, s := util.ProcessVersion(
			args[0],
			hotfixCmdOptions.VersionFile,
			hotfixCmdOptions.VersionFileType,
		)

		hotfix, err := glow.NewHotfix(version)
		util.ExitOnError(err)

		util.ExitOnError(g.Create(hotfix, rootCmdOptions.SkipChecks))
		util.ExitOnError(g.Checkout(hotfix))

		if util.IsSemanticVersion(args[0]) {
			util.ExitOnError(s.SetNextVersion(args[0]))
		} else {
			util.ExitOnError(s.SetVersion(version))
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
