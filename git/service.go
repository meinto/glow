package git

import "github.com/meinto/glow"

// Service describes all actions which can performed with git
type Service interface {
	CurrentBranch() (glow.Branch, error)
	Create(b glow.Branch) error
	Checkout(b glow.Branch) error
}

type service struct {
	Service
}

func NewGoGitService() Service {
	return service{goGitAdapter{}}
}

func NewNativeService(gitPath string) Service {
	return service{nativeGitAdapter{gitPath}}
}
