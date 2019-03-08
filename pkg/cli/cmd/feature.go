package cmd

import (
	"github.com/meinto/glow"

	"github.com/meinto/glow/git"
	"github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(featureCmd)
}

var featureCmd = &cobra.Command{
	Use:   "feature",
	Short: "create a feature branch",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		featureName := args[0]

		feature, err := glow.NewFeature(viper.GetString("author"), featureName)
		util.CheckForError(err, "NewFeature")

		g := git.NewGit()

		err = g.Create(feature)
		util.CheckForError(err, "Create")

		g.Checkout(feature)
		util.CheckForError(err, "Checkout")

		// r, err := git.PlainOpen(".")
		// util.CheckForError(err, "PlainOpen")

		// headRef, err := r.Head()
		// util.CheckForError(err, "Head")

		// refName := string(headRef.Name())
		// if !strings.Contains(refName, "develop") {
		// 	log.Println("You are not on the develop branch.")
		// 	log.Fatalf("Please switch branch...")
		// }

		// branchName := fmt.Sprintf("refs/heads/feature/%s/%s", viper.GetString("author"), feature)
		// ref := plumbing.NewHashReference(plumbing.ReferenceName(branchName), headRef.Hash())

		// err = r.Storer.SetReference(ref)
		// util.CheckForError(err, "SetReference")

		// w, err := r.Worktree()
		// util.CheckForError(err, "Worktree")

		// err = util.Checkout(w, branchName, util.ShouldUseNativeGitBinding("checkout"))
		// util.CheckForError(err, "Checkout")
	},
}
