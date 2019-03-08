package git

import (
	"os/exec"

	"github.com/meinto/glow"
	"github.com/pkg/errors"
)

// NativeGitAdapter implemented with native git
type NativeGitAdapter struct {
	gitPath string
}

// CurrentBranch returns the current branch name
func (g NativeGitAdapter) CurrentBranch() (glow.Branch, error) {
	return glow.Branch{}, errors.New("not implemented yet")
}

// Create a new branch
func (g NativeGitAdapter) Create(b glow.IBranch) error {
	return errors.New("not implemented yet")
}

// Checkout a branch
func (g NativeGitAdapter) Checkout(b glow.IBranch) error {
	cmd := exec.Command(g.gitPath, "checkout", b.ShortBranchName())
	err := cmd.Run()
	return errors.Wrap(err, "native checkout")
}
