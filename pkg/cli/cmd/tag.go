package cmd

import (
	"github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/meinto/glow/semver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(tagCmd)
}

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "create a tag of current version",
	Run: func(cmd *cobra.Command, args []string) {
		g, err := util.GetGitClient()
		util.CheckForError(err, "GetGitClient")

		pathToRepo, err := g.GitRepoPath()
		util.CheckForError(err, "semver GitRepoPath")

		s := semver.NewSemverService(
			pathToRepo,
			"/bin/bash",
			releaseCmdOptions.VersionFile,
			releaseCmdOptions.VersionFileType,
		)

		err = s.TagCurrentVersion()
		util.CheckForError(err, "TagCurrentVersion")
	},
}
