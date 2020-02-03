package semver

import (
	l "github.com/meinto/glow/logging"
)

type loggingService struct {
	next Service
}

func NewLoggingService(s Service) Service {
	l.Log().Trace(l.Fields{"service": s})
	return &loggingService{s}
}

func (s *loggingService) GetCurrentVersion() (version string, err error) {
	defer func() {
		l.Log().
			Info(l.Fields{"version": version}).
			Error(err)
	}()
	return s.next.GetCurrentVersion()
}

func (s *loggingService) GetNextVersion(versionType string) (version string, err error) {
	defer func() {
		l.Log().
			Info(l.Fields{
				"versionType": versionType,
				"version":     version,
			}).
			Error(err)
	}()
	return s.next.GetNextVersion(versionType)
}

func (s *loggingService) SetNextVersion(versionType string) (err error) {
	defer func() {
		l.Log().
			Info(l.Fields{"versionType": versionType}).
			Error(err)
	}()
	return s.next.SetNextVersion(versionType)
}

func (s *loggingService) SetVersion(version string) (err error) {
	defer func() {
		l.Log().
			Info(l.Fields{"version": version}).
			Error(err)
	}()
	return s.next.SetVersion(version)
}

func (s *loggingService) TagCurrentVersion() (err error) {
	defer func() {
		l.Log().Error(err)
	}()
	return s.next.TagCurrentVersion()
}
