package cmd

import (
	"github.com/meinto/glow/pkg/cli/cmd/internal/util"
	"github.com/meinto/glow/semver"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(tagCmd)
}

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "create a tag of current version",
	Run: func(cmd *cobra.Command, args []string) {
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

		util.ExitOnError(s.TagCurrentVersion())
	},
}
