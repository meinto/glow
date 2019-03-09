package githost

import (
	"github.com/go-kit/kit/log"
	"github.com/meinto/glow"
	"github.com/pkg/errors"
)

type loggingService struct {
	logger log.Logger
	next   Service
}

func NewLoggingService(l log.Logger, s Service) Service {
	return &loggingService{l, s}
}

func (s *loggingService) Close(b glow.Branch) error {
	return errors.New("not implemented yet")
}

func (s *loggingService) Publish(b glow.Branch) error {
	return errors.New("not implemented yet")
}
