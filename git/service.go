package git

import (
	"github.com/meinto/glow"
	"github.com/pkg/errors"
)

// Service describes all actions which can performed with git
type Service interface {
	SetCICDOrigin(origin string) error
	GitRepoPath() (string, error)
	CurrentBranch() (glow.Branch, error)
	BranchList() ([]glow.Branch, error)
	Fetch() error
	Push(setUpstream bool) error
	Create(b glow.Branch) error
	Checkout(b glow.Branch) error
	CleanupBranches(cleanupGone, cleanupUntracked bool) error
	CleanupTags(cleanupUntracked bool) error
}

type service struct {
	Service
}

func NewGoGitService() Service {
	return service{goGitAdapter{}}
}

func NewNativeService(gitPath string) (Service, error) {
	if gitPath == "" {
		return nil, errors.New("path to git not defined")
	}
	return service{nativeGitAdapter{gitPath}}, nil
}
