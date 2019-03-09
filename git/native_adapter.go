package git

import (
	"os/exec"

	"github.com/meinto/glow"
	"github.com/pkg/errors"
)

// NativeGitAdapter implemented with native git
type nativeGitAdapter struct {
	gitPath string
}

// CurrentBranch returns the current branch name
func (a nativeGitAdapter) CurrentBranch() (glow.Branch, error) {
	return glow.Branch{}, errors.New("not implemented yet")
}

// Create a new branch
func (a nativeGitAdapter) Create(b glow.IBranch) error {
	return errors.New("not implemented yet")
}

// Checkout a branch
func (a nativeGitAdapter) Checkout(b glow.IBranch) error {
	cmd := exec.Command(a.gitPath, "checkout", b.ShortBranchName())
	err := cmd.Run()
	return errors.Wrap(err, "native checkout")
}
