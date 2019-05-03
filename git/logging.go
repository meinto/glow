package git

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/meinto/glow"
)

type loggingService struct {
	logger log.Logger
	next   Service
}

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

// SetCICDOrigin for pipeline
func (s loggingService) SetCICDOrigin(origin string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "SetCICDOrigin", "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.SetCICDOrigin(origin)
}

// GitRepoPath returns the path to the root with the .git folder
func (s loggingService) GitRepoPath() (_ string, err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "GitRepoPath", "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.GitRepoPath()
}

// CurrentBranch returns the current branch name
func (s loggingService) CurrentBranch() (_ glow.Branch, err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "CurrentBranch", "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.CurrentBranch()
}

// BranchList returns a list of avalilable branches
func (s loggingService) BranchList() (_ []glow.Branch, err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "BranchList", "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.BranchList()
}

// Fetch changes
func (s loggingService) Fetch() (err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "Fetch", "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.Fetch()
}

// Add all changes
func (s loggingService) AddAll() (err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "AddAll", "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.AddAll()
}

// Commit added changes
func (s loggingService) Commit(message string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "Commit", "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.Commit(message)
}

// Push changes
func (s loggingService) Push(setUpstream bool) (err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "Push", "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.Push(setUpstream)
}

// Create a new branch
func (s loggingService) Create(b glow.Branch, skipChecks bool) (err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "Create", "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.Create(b, skipChecks)
}

// Checkout a branch
func (s loggingService) Checkout(b glow.Branch) (err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "Checkout", "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.Checkout(b)
}

// CleanupBranches removes all unused branches
func (s loggingService) CleanupBranches(cleanupGone, cleanupUntracked bool) (err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "CleanupBranches", "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.CleanupBranches(cleanupGone, cleanupUntracked)
}

// CleanupTags removes tags from local repo
func (s loggingService) CleanupTags(cleanupUntracked bool) (err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "CleanupTags", "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.CleanupTags(cleanupUntracked)
}
