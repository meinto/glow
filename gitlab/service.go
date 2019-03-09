package gitlab

import (
	"github.com/meinto/glow"
	"github.com/pkg/errors"
)

// Service describes all actions which can performed
// with the git hosting git service (gitlab etc.)
type Service interface {
	Close(b glow.IBranch) error
	Publish(b glow.IBranch) error
}

type service struct {
	endpoint string
	token    string
}

func (s *service) Close(b glow.IBranch) error {
	return errors.New("not implemented yet")
}

func (s *service) Publish(b glow.IBranch) error {
	return errors.New("not implemented yet")
}
