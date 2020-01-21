package cmd

import (
	. "github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/meinto/glow/semver"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(tagCmd)
}

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "create a tag of current version",
	Run: func(cmd *cobra.Command, args []string) {
		g, err := GetGitClient()
		ExitOnError(err)

		pathToRepo, _, _, err := g.GitRepoPath()
		ExitOnError(err)

		s := semver.NewSemverService(
			pathToRepo,
			"/bin/bash",
			releaseCmdOptions.VersionFile,
			releaseCmdOptions.VersionFileType,
		)
		s = semver.NewLoggingService(s)

		ExitOnError(s.TagCurrentVersion())
	},
}
