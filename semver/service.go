package semver

import (
	"bytes"
	"os/exec"

	"github.com/pkg/errors"
)

type Service interface {
	GetNextVersion(versionType string) (string, error)
	SetNextVersion(versionType string) error
}

type service struct {
	pathToGit    string
	pathToSemver string
}

func NewSemverService(pathToSemver string) Service {
	return &service{
		pathToSemver: pathToSemver,
	}
}

func NewGitSemverService(pathToGit string) Service {
	return &service{
		pathToGit: pathToGit,
	}
}

func (s *service) GetNextVersion(versionType string) (string, error) {
	var cmd *exec.Cmd
	if s.pathToGit != "" {
		cmd = exec.Command(s.pathToGit, "semver", "get", versionType, "-r")
	}
	if s.pathToSemver != "" {
		cmd = exec.Command(s.pathToSemver, "get", versionType, "-r")
	}

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	return stdout.String(), errors.Wrap(err, "Error getting next version; Is semver installed?")
}

func (s *service) SetNextVersion(versionType string) error {
	var cmd *exec.Cmd
	if s.pathToGit != "" {
		cmd = exec.Command(s.pathToGit, "semver", "version", versionType)
	}
	if s.pathToSemver != "" {
		cmd = exec.Command(s.pathToSemver, "version", versionType)
	}
	err := cmd.Run()
	return errors.Wrap(err, "Error getting setting version; Is semver installed?")
}
