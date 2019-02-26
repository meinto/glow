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
	rootCmd.AddCommand(hotfixCmd)
}

var hotfixCmd = &cobra.Command{
	Use:   "hotfix",
	Short: "create a hotfix branch",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hotix := args[0]

		r, err := git.PlainOpen(".")
		util.CheckForError(err, "PlainOpen")

		headRef, err := r.Head()
		util.CheckForError(err, "Head")

		refName := string(headRef.Name())
		if !strings.Contains(refName, "master") {
			log.Println("You are not on the master branch.")
			log.Fatalf("Please switch branch...")
		}

		branchName := fmt.Sprintf("refs/heads/hotfix/%s/%s", viper.GetString("author"), hotix)
		ref := plumbing.NewHashReference(plumbing.ReferenceName(branchName), headRef.Hash())

		err = r.Storer.SetReference(ref)
		util.CheckForError(err, "SetReference")

		w, err := r.Worktree()
		util.CheckForError(err, "Worktree")

		err = w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName(branchName),
		})
		util.CheckForError(err, "Checkout")
	},
}
