package gitprovider

import (
	"github.com/meinto/glow"
	l "github.com/meinto/glow/logging"
)

type loggingService struct {
	next Service
}

func NewLoggingService(s Service) Service {
	defer func() {
		l.Log().Info(l.Fields{"service": s})
	}()
	return &loggingService{s}
}

func (s *loggingService) Close(b glow.Branch) (err error) {
	defer func() {
		l.Log().
			Info(l.Fields{"branch": b.BranchName()}).
			Error(err)
	}()
	return s.next.Close(b)
}

func (s *loggingService) Publish(b glow.Branch) (err error) {
	defer func() {
		l.Log().
			Info(l.Fields{"branch": b.BranchName()}).
			Error(err)
	}()
	return s.next.Publish(b)
}

func (s *loggingService) GetCIBranch() (branch glow.Branch) {
	defer func() {
		l.Log().Info(l.Fields{"branch": branch.BranchName()})
	}()
	return s.next.GetCIBranch()
}

func (s *loggingService) DetectCICDOrigin() (cicdOrigin string, err error) {
	defer func() {
		l.Log().
			Info(l.Fields{"cicdOrigin": cicdOrigin}).
			Error(err)
	}()
	return s.next.DetectCICDOrigin()
}
