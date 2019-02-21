package cmd

import (
	"fmt"

	"github.com/meinto/glow/cmd/utils"
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
	Run: func(cmd *cobra.Command, args []string) {
		feature := args[0]

		r, err := git.PlainOpen(".")
		utils.CheckForError(err, "PlainOpen")

		headRef, err := r.Head()
		utils.CheckForError(err, "Head")

		branchName := fmt.Sprintf("refs/heads/features/%s/%s", viper.GetString("author"), feature)
		ref := plumbing.NewHashReference(plumbing.ReferenceName(branchName), headRef.Hash())

		err = r.Storer.SetReference(ref)
		utils.CheckForError(err, "SetReference")

		w, err := r.Worktree()
		utils.CheckForError(err, "Worktree")

		err = w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName(branchName),
		})
		utils.CheckForError(err, "Checkout")
	},
}
