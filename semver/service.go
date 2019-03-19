package semver

import (
	"bytes"
	"os/exec"

	semver "github.com/meinto/git-semver"
	"github.com/meinto/git-semver/file"
	"github.com/pkg/errors"
)

type Service interface {
	GetNextVersion(versionType string) (string, error)
	SetNextVersion(versionType string) error
}

type service struct {
	pathToRepo      string
	pathToGit       string
	versionFile     string
	versionFileType string
}

func NewSemverService(pathToRepo, pathToGit string) Service {
	return &service{
		pathToRepo: pathToRepo,
		pathToGit:  pathToGit,
	}
}

func (s *service) GetNextVersion(versionType string) (string, error) {
	versionFilepath := s.pathToRepo + "/" + s.versionFile
	fs := file.NewVersionFileService(versionFilepath)

	currentVersion, err := fs.ReadVersionFromFile(s.versionFileType)
	if err != nil {
		return "", err
	}

	vs, err := semver.NewVersion(currentVersion)
	if err != nil {
		return "", err
	}

	nextVersion, err := vs.Get(versionType)
	return nextVersion, errors.Wrap(err, "GetNextVersion")
}

func (s *service) SetNextVersion(versionType string) error {
	var cmd *exec.Cmd
	if s.pathToGit != "" {
		cmd = exec.Command(s.pathToGit, "semver", "version", versionType)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return errors.Wrap(err, stderr.String())
}
