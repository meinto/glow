package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/meinto/glow/cmd/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func init() {
	rootCmd.AddCommand(fixCmd)
}

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "create a fix branch",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fix := args[0]

		r, err := git.PlainOpen(".")
		util.CheckForError(err, "PlainOpen")

		headRef, err := r.Head()
		util.CheckForError(err, "Head")

		refName := string(headRef.Name())
		if !strings.Contains(refName, "release/") {
			log.Println("You are not on a release branch.")
			log.Fatalf("Please switch branch...")
		}

		branchName := fmt.Sprintf("refs/heads/fix/%s/%s", viper.GetString("author"), fix)
		ref := plumbing.NewHashReference(plumbing.ReferenceName(branchName), headRef.Hash())

		err = r.Storer.SetReference(ref)
		util.CheckForError(err, "SetReference")

		w, err := r.Worktree()
		util.CheckForError(err, "Worktree")

		err = util.Checkout(w, branchName, util.ShouldUseNativeGitBinding("checkout"))
		util.CheckForError(err, "Checkout")
	},
}
