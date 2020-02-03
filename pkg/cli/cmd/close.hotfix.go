package cmd

import (
	l "github.com/meinto/glow/logging"

	"github.com/meinto/glow"
	"github.com/meinto/glow/pkg/cli/cmd/internal/util"
	"github.com/meinto/glow/semver"
	"github.com/spf13/cobra"
)

func init() {
	closeCmd.Cmd().AddCommand(closeHotfixCmd)
	util.AddFlagsForMergeRequests(closeHotfixCmd)
}

var closeHotfixCmd = &cobra.Command{
	Use:   "hotfix",
	Short: "close a release branch",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]

		if version == "current" {
			g, err := util.GetGitClient()
			util.ExitOnError(err)

			pathToRepo, _, _, err := g.GitRepoPath()
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

		gp, err := util.GetGitProvider()
		util.ExitOnError(err)

		currentBranch, err := glow.NewHotfix(version)
		util.ExitOnError(err)

		err = gp.Close(currentBranch)
		if !util.MergeRequestFlags.Gracefully {
			util.ExitOnError(err)
		} else {
			l.Log().Error(err)
		}
	},
}
