package git

import (
	"github.com/meinto/glow"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// GoGitAdapter implemented with go-git
type goGitAdapter struct{}

// CurrentBranch returns the current branch name
func (a goGitAdapter) CurrentBranch() (glow.Branch, error) {
	r, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return glow.Branch{}, errors.Wrap(err, "error opening repository")
	}

	headRef, err := r.Head()
	if err != nil {
		return glow.Branch{}, errors.Wrap(err, "error getting current branch")
	}

	refName := string(headRef.Name())
	return glow.NewBranch(refName)
}

// Create a new branch
func (a goGitAdapter) Create(b glow.IBranch) error {
	r, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return errors.Wrap(err, "error opening repository")
	}

	headRef, err := r.Head()
	if err != nil {
		return errors.Wrap(err, "error getting current branch")
	}

	refName := string(headRef.Name())
	if !b.CreationIsAllowedFrom(refName) {
		return errors.New("You are not on the develop branch.\nPlease switch branch...\n")
	}

	ref := plumbing.NewHashReference(plumbing.ReferenceName(b.BranchName()), headRef.Hash())

	err = r.Storer.SetReference(ref)
	return errors.Wrap(err, "error while creating branch")
}

// Checkout a branch
func (a goGitAdapter) Checkout(b glow.IBranch) error {
	r, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return errors.Wrap(err, "error opening repository")
	}

	w, err := r.Worktree()
	if err != nil {
		return errors.Wrap(err, "error getting worktree")
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(b.BranchName()),
	})
	return errors.Wrapf(err, "error while checkout branch %s", b.BranchName())
}
