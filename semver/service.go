package semver

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

type Service interface {
	GetNextVersion(versionType string) (string, error)
	SetNextVersion(versionType string) error
}

type service struct {
	pathToRepo   string
	pathToGit    string
	pathToSemver string
}

func NewSemverService(pathToRepo, pathToSemver string) Service {
	return &service{
		pathToRepo:   pathToRepo,
		pathToSemver: pathToSemver,
	}
}

func NewGitSemverService(pathToRepo, pathToGit string) Service {
	return &service{
		pathToRepo: pathToRepo,
		pathToGit:  pathToGit,
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

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return strings.TrimSuffix(stdout.String(), "\n"), errors.Wrap(err, stderr.String())
}

func (s *service) SetNextVersion(versionType string) error {
	var cmd *exec.Cmd
	if s.pathToGit != "" {
		cmd = exec.Command(s.pathToGit, "semver", "version", versionType)
	}
	if s.pathToSemver != "" {
		cmd = exec.Command(s.pathToSemver, "version", versionType)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return errors.Wrap(err, stderr.String())
}
