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

// CurrentBranch returns the current branch name
func (s loggingService) CurrentBranch() (b glow.Branch, err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "CurrentBranch", "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.CurrentBranch()
}

// BranchList returns a list of avalilable branches
func (s loggingService) BranchList() (b []glow.Branch, err error) {
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

// Create a new branch
func (s loggingService) Create(b glow.Branch) (err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "Create", "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.Create(b)
}

// Checkout a branch
func (s loggingService) Checkout(b glow.Branch) (err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "Checkout", "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.Checkout(b)
}
