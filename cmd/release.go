package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var postReleaseScript string

func init() {
	rootCmd.AddCommand(releaseCmd)
	releaseCmd.Flags().StringVar(&postReleaseScript, "postRelease", "", "script that executes after switching to release branch")
}

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "create a release branch",
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0] // should be semver

		r, err := git.PlainOpen(".")
		CheckForError(err, "PlainOpen")

		headRef, err := r.Head()
		CheckForError(err, "Head")

		branchName := fmt.Sprintf("refs/heads/release/v%s", version)
		ref := plumbing.NewHashReference(plumbing.ReferenceName(branchName), headRef.Hash())

		err = r.Storer.SetReference(ref)
		CheckForError(err, "SetReference")

		w, err := r.Worktree()
		CheckForError(err, "Worktree")

		err = w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName(branchName),
		})
		CheckForError(err, "Checkout")

		if postReleaseScript != "" {
			postRelease(version)
		}
	},
}

func postRelease(version string) {
	pathToFile, err := filepath.Abs(postReleaseScript)
	if err != nil {
		log.Println("cannot find post-release script", err)
	}
	cmd := exec.Command(pathToFile, version)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Println("error while executing post-release script", err)
	}
	log.Println("post release:")
	log.Println(out.String())
}
