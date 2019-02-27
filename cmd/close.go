package cmd

import (
	"log"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/meinto/glow/cmd/util"
	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func init() {
	rootCmd.AddCommand(finishCmd)
	util.AddFlagsForMergeRequests(finishCmd)
}

var finishCmd = &cobra.Command{
	Use:   "close",
	Short: "close a branch",
	Run: func(cmd *cobra.Command, args []string) {

		r, err := git.PlainOpen(".")
		util.CheckForError(err, "PlainOpen")

		r.Fetch(&git.FetchOptions{})

		headRef, err := r.Head()
		refName := string(headRef.Name())

		if strings.Contains(refName, "feature/") ||
			strings.Contains(refName, "release/") ||
			strings.Contains(refName, "hotix/") {
			util.CreateMergeRequest(refName, "develop")
		}

		if strings.Contains(refName, "fix/") {
			var releaseBranches = make([]string, 0)
			iterator, _ := r.References()
			iterator.ForEach(func(ref *plumbing.Reference) error {
				name := string(ref.Name())
				if strings.Contains(name, "refs/heads/release/v") {
					releaseBranches = append(releaseBranches, string(ref.Name()))
				}
				return nil
			})
			if len(releaseBranches) > 0 {
				index, err := mergeWith(releaseBranches)
				if err != nil {
					log.Fatalf(err.Error())
				}
				mergeWith := releaseBranches[index]
				util.CreateMergeRequest(refName, mergeWith)
			} else {
				log.Println("There is no release branch you could merge with.")
			}
		}
	},
}

func mergeWith(refs []string) (int, error) {
	prompt := promptui.Select{
		Label: "Which branch do you want for merge?",
		Items: refs,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return -1, err
	}

	return index, nil
}
