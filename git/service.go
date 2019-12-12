package git

import (
	"github.com/meinto/glow"
	"github.com/meinto/glow/cmd"
)

// Service describes all actions which can performed with git
type Service interface {
	SetCICDOrigin(origin string) error
	GitRepoPath() (string, error)
	CurrentBranch() (glow.Branch, error)
	BranchList() ([]glow.Branch, error)
	Fetch() error
	AddAll() error
	Commit(message string) error
	Push(setUpstream bool) error
	Create(b glow.Branch, skipChecks bool) error
	Checkout(b glow.Branch) error
	CleanupBranches(cleanupGone, cleanupUntracked bool) error
	CleanupTags(cleanupUntracked bool) error
	RemoteBranchExists(branchName string) error
}

type service struct {
	Service
}

func NewGoGitService() Service {
	return service{goGitAdapter{}}
}

func NewNativeService(cmdExecutor cmd.CmdExecutor) Service {
	return service{nativeGitAdapter{cmdExecutor}}
}
