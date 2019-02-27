package cmd

import (
	"strings"

	"github.com/meinto/glow/cmd/util"
	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
)

func init() {
	rootCmd.AddCommand(publishCmd)
	util.AddFlagsForMergeRequests(publishCmd)
}

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish a release branch",
	Run: func(cmd *cobra.Command, args []string) {

		r, err := git.PlainOpen(".")
		util.CheckForError(err, "PlainOpen")

		r.Fetch(&git.FetchOptions{})

		headRef, err := r.Head()
		refName := string(headRef.Name())

		if strings.Contains(refName, "release/") ||
			strings.Contains(refName, "hotfix/") {
			util.CreateMergeRequest(refName, "master")
		}
	},
}
