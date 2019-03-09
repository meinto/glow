package githost

import (
	"github.com/meinto/glow"
	"github.com/pkg/errors"
)

type gitlabAdapter struct {
	service
}

func (s *gitlabAdapter) Close(b glow.IBranch) error {
	return errors.New("not implemented yet")
}

func (s *gitlabAdapter) Publish(b glow.IBranch) error {
	return errors.New("not implemented yet")
}
