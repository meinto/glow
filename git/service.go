package git

import (
	"bytes"

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
	Stash() error
	StashPop() (stdout, stderr bytes.Buffer, err error)
	Commit(message string) error
	Push(setUpstream bool) (stdout, stderr bytes.Buffer, err error)
	Create(b glow.Branch, skipChecks bool) error
	Checkout(b glow.Branch) error
	CleanupBranches(cleanupGone, cleanupUntracked bool) error
	CleanupTags(cleanupUntracked bool) error
	RemoteBranchExists(branchName string) error
}

type NativeService interface {
	Service
	CMDExecutor() cmd.CmdExecutor
}

type service struct{ Service }

func NewGoGitService() Service {
	return service{goGitAdapter{}}
}

func NewNativeService(cmdExecutor cmd.CmdExecutor) Service {
	return service{nativeGitAdapter{cmdExecutor}}
}
