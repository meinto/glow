package gitprovider

import (
	"github.com/meinto/glow"
	l "github.com/meinto/glow/logging"
	"github.com/sirupsen/logrus"
)

type loggingService struct {
	next Service
}

func NewLoggingService(s Service) Service {
	return &loggingService{s}
}

func (s *loggingService) Close(b glow.Branch) (err error) {
	defer func() {
		l.Log().WithFields(logrus.Fields{
			"branch": b.BranchName(),
			"error":  err,
		}).Info()
	}()
	return s.next.Close(b)
}

func (s *loggingService) Publish(b glow.Branch) (err error) {
	defer func() {
		l.Log().WithFields(logrus.Fields{
			"branch": b.BranchName(),
			"error":  err,
		}).Info()
	}()
	return s.next.Publish(b)
}

func (s *loggingService) GetCIBranch() (branch glow.Branch, err error) {
	defer func() {
		l.Log().WithFields(logrus.Fields{
			"branch": branch.BranchName(),
			"error":  err,
		}).Info()
	}()
	return s.next.GetCIBranch()
}

func (s *loggingService) DetectCICDOrigin() (cicdOrigin string, err error) {
	defer func() {
		l.Log().WithFields(logrus.Fields{
			"cicdOrigin": cicdOrigin,
			"error":      err,
		}).Info()
	}()
	return s.next.DetectCICDOrigin()
}
