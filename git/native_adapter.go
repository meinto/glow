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

// GitRepoPath returns the path to the root with the .git folder
func (a nativeGitAdapter) GitRepoPath() (string, error) {
	return "", errors.New("not implemented yet")
}

// CurrentBranch returns the current branch name
func (a nativeGitAdapter) CurrentBranch() (glow.Branch, error) {
	return nil, errors.New("not implemented yet")
}

// BranchList returns a list of avalilable branches
func (a nativeGitAdapter) BranchList() ([]glow.Branch, error) {
	return nil, errors.New("not implemented yet")
}

// Fetch changes
func (a nativeGitAdapter) Fetch() error {
	cmd := exec.Command(a.gitPath, "fetch")
	err := cmd.Run()
	return errors.Wrap(err, "native Fetch")
}

// Create a new branch
func (a nativeGitAdapter) Create(b glow.Branch) error {
	cmd := exec.Command(a.gitPath, "branch", b.ShortBranchName())
	err := cmd.Run()
	return errors.Wrap(err, "native Create")
}

// Checkout a branch
func (a nativeGitAdapter) Checkout(b glow.Branch) error {
	cmd := exec.Command(a.gitPath, "checkout", b.ShortBranchName())
	err := cmd.Run()
	return errors.Wrap(err, "native Checkout")
}
