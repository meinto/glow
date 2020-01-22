package git

import (
	"github.com/meinto/glow"
	l "github.com/meinto/glow/logging"
)

type loggingService struct {
	next Service
}

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(s Service) Service {
	defer func() {
		l.Log().Info(l.Fields{"service": s})
	}()
	return &loggingService{s}
}

// SetCICDOrigin for pipeline
func (s loggingService) SetCICDOrigin(origin string) (stdout, stderr string, err error) {
	defer func() {
		l.Log().
			StdoutFields(stdout, l.Fields{"origin": origin}).
			Stderr(stderr, err)
	}()
	return s.next.SetCICDOrigin(origin)
}

// GitRepoPath returns the path to the root with the .git folder
func (s loggingService) GitRepoPath() (repoPath string, stdout, stderr string, err error) {
	defer func() {
		l.Log().
			StdoutFields(stdout, l.Fields{"repoPath": repoPath}).
			Stderr(stderr, err)
	}()
	return s.next.GitRepoPath()
}

// CurrentBranch returns the current branch name
func (s loggingService) CurrentBranch() (branch glow.Branch, stdout, stderr string, err error) {
	defer func() {
		l.Log().
			StdoutFields(stdout, l.Fields{"branchName": branch.BranchName()}).
			Stderr(stderr, err)
	}()
	return s.next.CurrentBranch()
}

// BranchList returns a list of avalilable branches
func (s loggingService) BranchList() (_ []glow.Branch, stdout, stderr string, err error) {
	defer func() {
		l.Log().
			Stdout(stdout).
			Stderr(stderr, err)
	}()
	return s.next.BranchList()
}

// Fetch changes
func (s loggingService) Fetch() (stdout, stderr string, err error) {
	defer func() {
		l.Log().
			Stdout(stdout).
			Stderr(stderr, err)
	}()
	return s.next.Fetch()
}

// Add all changes
func (s loggingService) AddAll() (stdout, stderr string, err error) {
	defer func() {
		l.Log().
			Stdout(stdout).
			Stderr(stderr, err)
	}()
	return s.next.AddAll()
}

// Stash all changes
func (s loggingService) Stash() (stdout, stderr string, err error) {
	defer func() {
		l.Log().
			Stdout(stdout).
			Stderr(stderr, err)
	}()
	return s.next.Stash()
}

// Pop all stashed changes
func (s loggingService) StashPop() (stdout, stderr string, err error) {
	defer func() {
		l.Log().
			Stdout(stdout).
			Stderr(stderr, err)
	}()
	return s.next.StashPop()
}

// Commit added changes
func (s loggingService) Commit(message string) (stdout, stderr string, err error) {
	defer func() {
		l.Log().
			StdoutFields(stdout, l.Fields{"message": message}).
			Stderr(stderr, err)
	}()
	return s.next.Commit(message)
}

// Push changes
func (s loggingService) Push(setUpstream bool) (stdout, stderr string, err error) {
	defer func() {
		l.Log().
			StdoutFields(stdout, l.Fields{"setUpstream": setUpstream}).
			Stderr(stderr, err)
	}()
	return s.next.Push(setUpstream)
}

// Create a new branch
func (s loggingService) Create(b glow.Branch, skipChecks bool) (stdout, stderr string, err error) {
	defer func() {
		l.Log().
			StdoutFields(stdout, l.Fields{
				"branchName": b.BranchName(),
				"skipChecks": skipChecks,
			}).
			Stderr(stderr, err)
	}()
	return s.next.Create(b, skipChecks)
}

// Checkout a branch
func (s loggingService) Checkout(b glow.Branch) (stdout, stderr string, err error) {
	defer func() {
		l.Log().
			StdoutFields(stdout, l.Fields{"branchName": b.BranchName()}).
			Stderr(stderr, err)
	}()
	return s.next.Checkout(b)
}

// CleanupBranches removes all unused branches
func (s loggingService) CleanupBranches(cleanupGone, cleanupUntracked bool) (stdout, stderr string, err error) {
	defer func() {
		l.Log().
			StdoutFields(stdout, l.Fields{
				"cleanupGone":      cleanupGone,
				"cleanupUntracked": cleanupUntracked,
			}).
			Stderr(stderr, err)
	}()
	return s.next.CleanupBranches(cleanupGone, cleanupUntracked)
}

// CleanupTags removes tags from local repo
func (s loggingService) CleanupTags(cleanupUntracked bool) (stdout, stderr string, err error) {
	defer func() {
		l.Log().
			StdoutFields(stdout, l.Fields{"cleanupUntracked": cleanupUntracked}).
			Stderr(stderr, err)
	}()
	return s.next.CleanupTags(cleanupUntracked)
}

func (s loggingService) RemoteBranchExists(branchName string) (stdout, stderr string, err error) {
	defer func() {
		l.Log().
			StdoutFields(stdout, l.Fields{"branchName": branchName}).
			Stderr(stderr, err)
	}()
	return s.next.RemoteBranchExists(branchName)
}
