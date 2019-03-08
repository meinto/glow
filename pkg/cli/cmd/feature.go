package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func init() {
	rootCmd.AddCommand(featureCmd)
}

var featureCmd = &cobra.Command{
	Use:   "feature",
	Short: "create a feature branch",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		feature := args[0]

		r, err := git.PlainOpen(".")
		util.CheckForError(err, "PlainOpen")

		headRef, err := r.Head()
		util.CheckForError(err, "Head")

		refName := string(headRef.Name())
		if !strings.Contains(refName, "develop") {
			log.Println("You are not on the develop branch.")
			log.Fatalf("Please switch branch...")
		}

		branchName := fmt.Sprintf("refs/heads/feature/%s/%s", viper.GetString("author"), feature)
		ref := plumbing.NewHashReference(plumbing.ReferenceName(branchName), headRef.Hash())

		err = r.Storer.SetReference(ref)
		util.CheckForError(err, "SetReference")

		w, err := r.Worktree()
		util.CheckForError(err, "Worktree")

		err = util.Checkout(w, branchName, util.ShouldUseNativeGitBinding("checkout"))
		util.CheckForError(err, "Checkout")
	},
}
