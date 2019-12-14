package git

import (
	"bytes"

	"github.com/meinto/glow"
	"github.com/meinto/glow/cmd"
)

// Service describes all actions which can performed with git
type Service interface {
	SetCICDOrigin(origin string) (stdout, stderr bytes.Buffer, err error)
	GitRepoPath() (path string, stdout, stderr bytes.Buffer, err error)
	CurrentBranch() (branch glow.Branch, stdout, stderr bytes.Buffer, err error)
	BranchList() (branchList []glow.Branch, stdout, stderr bytes.Buffer, err error)
	Fetch() (stdout, stderr bytes.Buffer, err error)
	AddAll() (stdout, stderr bytes.Buffer, err error)
	Stash() (stdout, stderr bytes.Buffer, err error)
	StashPop() (stdout, stderr bytes.Buffer, err error)
	Commit(message string) (stdout, stderr bytes.Buffer, err error)
	Push(setUpstream bool) (stdout, stderr bytes.Buffer, err error)
	Create(b glow.Branch, skipChecks bool) (stdout, stderr bytes.Buffer, err error)
	Checkout(b glow.Branch) (stdout, stderr bytes.Buffer, err error)
	CleanupBranches(cleanupGone, cleanupUntracked bool) (stdout, stderr bytes.Buffer, err error)
	CleanupTags(cleanupUntracked bool) (stdout, stderr bytes.Buffer, err error)
	RemoteBranchExists(branchName string) (stdout, stderr bytes.Buffer, err error)
}

type NativeService interface {
	Service
	CMDExecutor() cmd.CmdExecutor
}

type service struct{ Service }

func NewNativeService(cmdExecutor cmd.CmdExecutor) Service {
	return service{nativeGitAdapter{cmdExecutor}}
}
