package util

import (
	"os/exec"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func trimBranchPrefix(branchName string) string {
	return strings.TrimPrefix(branchName, "refs/heads/")
}

func Checkout(w *git.Worktree, branchName string, useNativeCli bool) error {
	if useNativeCli {
		cmd := exec.Command(viper.GetString("gitPath"), "checkout", trimBranchPrefix(branchName))
		return cmd.Run()
	} else {
		return w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName(branchName),
		})
	}
}
