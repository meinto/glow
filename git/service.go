package git

import (
	"github.com/meinto/glow"
	"github.com/meinto/glow/cmd"
)

// Service describes all actions which can performed with git
type Service interface {
	SetCICDOrigin(origin string) (stdout, stderr string, err error)
	GitRepoPath() (path string, stdout, stderr string, err error)
	CurrentBranch() (branch glow.Branch, stdout, stderr string, err error)
	BranchList() (branchList []glow.Branch, stdout, stderr string, err error)
	Fetch() (stdout, stderr string, err error)
	AddAll() (stdout, stderr string, err error)
	Stash() (stdout, stderr string, err error)
	StashPop() (stdout, stderr string, err error)
	Commit(message string) (stdout, stderr string, err error)
	Push(setUpstream bool) (stdout, stderr string, err error)
	Create(b glow.Branch, skipChecks bool) (stdout, stderr string, err error)
	Checkout(b glow.Branch) (stdout, stderr string, err error)
	CleanupBranches(cleanupGone, cleanupUntracked bool) (stdout, stderr string, err error)
	CleanupTags(cleanupUntracked bool) (stdout, stderr string, err error)
	RemoteBranchExists(branchName string) (stdout, stderr string, err error)
}

type NativeService interface {
	Service
	CMDExecutor() cmd.CmdExecutor
}

type service struct{ Service }

func NewNativeService(cmdExecutor cmd.CmdExecutor) Service {
	return service{nativeGitAdapter{cmdExecutor}}
}
