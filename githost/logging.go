package githost

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/meinto/glow"
)

type loggingService struct {
	logger log.Logger
	next   Service
}

func NewLoggingService(l log.Logger, s Service) Service {
	return &loggingService{l, s}
}

func (s *loggingService) Close(b glow.Branch) (err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "Close", "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.Close(b)
}

func (s *loggingService) Publish(b glow.Branch) (err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "Publish", "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.Publish(b)
}
