package cmd

import (
	"github.com/meinto/glow"
	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
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

type HotfixCommand struct {
	command.Service
}

func (cmd *HotfixCommand) PostSetup(parent command.Service) command.Service {
	parent.Add(cmd)
	cmd.Cmd().Flags().BoolVar(&hotfixCmdOptions.Push, "push", false, "push created hotfix branch")
	cmd.Cmd().Flags().StringVar(&hotfixCmdOptions.PostHotfixScript, "postHotfix", "", "script that executes after switching to hotfix branch")
	cmd.Cmd().Flags().StringArrayVar(&hotfixCmdOptions.PostHotfixCommand, "postHotfixCommand", []string{}, "commands which should be executed after switching to hotfix branch")

	cmd.Cmd().Flags().StringVar(&hotfixCmdOptions.VersionFile, "versionFile", "VERSION", "name of git-semver version file")
	cmd.Cmd().Flags().StringVar(&hotfixCmdOptions.VersionFileType, "versionFileType", "raw", "git-semver version file type")
	return cmd
}

var hotfixCmd = SetupHotfixCommand(RootCmd)

func SetupHotfixCommand(parent command.Service) command.Service {
	return command.Setup(&HotfixCommand{
		&command.Command{
			Command: &cobra.Command{
				Use:   "hotfix",
				Short: "create a hotfix branch",
				Args:  cobra.MinimumNArgs(1),
			},
			Run: func(cmd command.Service, args []string) {
				pathToRepo, _, _, err := cmd.GitClient().GitRepoPath()
				util.ExitOnError(err)

				version, s := util.ProcessVersion(
					args[0],
					hotfixCmdOptions.VersionFile,
					hotfixCmdOptions.VersionFileType,
					pathToRepo,
				)

				hotfix, err := glow.NewHotfix(version)
				util.ExitOnError(err)

				util.ExitOnError(cmd.GitClient().Create(hotfix, rootCmdOptions.SkipChecks))
				util.ExitOnError(cmd.GitClient().Checkout(hotfix))

				if util.IsSemanticVersion(args[0]) {
					util.ExitOnError(s.SetNextVersion(args[0]))
				} else {
					util.ExitOnError(s.SetVersion(version))
				}
			},
			PostRun: func(cmd command.Service, args []string) {
				util.PostRunWithCurrentVersion(
					hotfixCmdOptions.VersionFile,
					hotfixCmdOptions.VersionFileType,
					hotfixCmdOptions.PostHotfixScript,
					hotfixCmdOptions.PostHotfixCommand,
					hotfixCmdOptions.Push,
				)
			},
		},
	}, parent)
}
