package semver

import (
	l "github.com/meinto/glow/logging"
	"github.com/sirupsen/logrus"
)

type loggingService struct {
	next Service
}

func NewLoggingService(s Service) Service {
	return &loggingService{s}
}

func (s *loggingService) GetCurrentVersion() (version string, err error) {
	defer func() {
		l.Log().WithFields(logrus.Fields{
			"version": version,
			"error":   err,
		}).Info()
	}()
	return s.next.GetCurrentVersion()
}

func (s *loggingService) GetNextVersion(versionType string) (version string, err error) {
	defer func() {
		l.Log().WithFields(logrus.Fields{
			"versionType": versionType,
			"version":     version,
			"error":       err,
		}).Info()
	}()
	return s.next.GetNextVersion(versionType)
}

func (s *loggingService) SetNextVersion(versionType string) (err error) {
	defer func() {
		l.Log().WithFields(logrus.Fields{
			"versionType": versionType,
			"error":       err,
		}).Info()
	}()
	return s.next.SetNextVersion(versionType)
}

func (s *loggingService) SetVersion(version string) (err error) {
	defer func() {
		l.Log().WithFields(logrus.Fields{
			"version": version,
			"error":   err,
		}).Info()
	}()
	return s.next.SetVersion(version)
}

func (s *loggingService) TagCurrentVersion() (err error) {
	defer func() {
		l.Log().WithFields(logrus.Fields{
			"error": err,
		}).Info()
	}()
	return s.next.TagCurrentVersion()
}
