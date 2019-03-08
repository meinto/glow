package git

import (
	"github.com/meinto/glow"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type GoGitAdapter struct{}

func (g GoGitAdapter) CurrentBranch() (glow.Branch, error) {
	r, err := git.PlainOpen(".")
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

func (g GoGitAdapter) Create(b glow.IBranch) error {
	r, err := git.PlainOpen(".")
	if err != nil {
		return errors.Wrap(err, "error opening repository")
	}

	headRef, err := r.Head()
	if err != nil {
		return errors.Wrap(err, "error getting current branch")
	}

	refName := string(headRef.Name())
	if b.CreationIsAllowedFrom(refName) {
		return errors.New("You are not on the develop branch.\nPlease switch branch...")
	}

	ref := plumbing.NewHashReference(plumbing.ReferenceName(b.BranchName()), headRef.Hash())

	err = r.Storer.SetReference(ref)
	return errors.Wrap(err, "error while creating branch")
}

func (g GoGitAdapter) Checkout(b glow.IBranch) error {
	return errors.New("not implemented yet")
}
