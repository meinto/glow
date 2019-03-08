package git

import (
	"github.com/meinto/glow"
	"github.com/pkg/errors"
)

// NativeGitAdapter implemented with native git
type NativeGitAdapter struct{}

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
	return errors.New("not implemented yet")
}
