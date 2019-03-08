package git

import (
	"github.com/meinto/glow"
	"github.com/pkg/errors"
)

type NativeGitAdapter struct{}

func (g NativeGitAdapter) CurrentBranch() (glow.Branch, error) {
	return glow.Branch{}, errors.New("not implemented yet")
}

func (g NativeGitAdapter) Create(b glow.IBranch) error {
	return errors.New("not implemented yet")
}

func (g NativeGitAdapter) Checkout(b glow.IBranch) error {
	return errors.New("not implemented yet")
}
