package semver

import (
	"log"
	"strings"

	semver "github.com/meinto/git-semver"
	"github.com/meinto/git-semver/file"
	"github.com/meinto/git-semver/git"
	"github.com/pkg/errors"
)

type Service interface {
	GetCurrentVersion() (string, error)
	GetNextVersion(versionType string) (string, error)
	SetNextVersion(versionType string) error
	TagCurrentVersion() error
}

type service struct {
	pathToRepo      string
	pathToGit       string
	versionFile     string
	versionFileType string
}

func NewSemverService(pathToRepo, pathToGit, versionFile, versionFileType string) Service {
	return &service{
		pathToRepo:      pathToRepo,
		pathToGit:       pathToGit,
		versionFile:     versionFile,
		versionFileType: versionFileType,
	}
}

func (s *service) GetCurrentVersion() (string, error) {
	versionFilepath := s.pathToRepo + "/" + s.versionFile
	fs := file.NewVersionFileService(versionFilepath)

	currentVersion, err := fs.ReadVersionFromFile(s.versionFileType)
	return strings.TrimSpace(currentVersion), errors.Wrap(err, "GetCurrentVersion")
}

func (s *service) GetNextVersion(versionType string) (string, error) {
	currentVersion, err := s.GetCurrentVersion()
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
	versionFilepath := s.pathToRepo + "/" + s.versionFile
	fs := file.NewVersionFileService(versionFilepath)

	currentVersion, err := s.GetCurrentVersion()
	if err != nil {
		return err
	}

	vs, err := semver.NewVersion(currentVersion)
	if err != nil {
		return err
	}

	nextVersion, err := vs.SetNext(versionType)
	if err != nil {
		return err
	}

	log.Println("new version will be: ", nextVersion)

	err = fs.WriteVersionFile(s.versionFileType, nextVersion)
	return err
}

func (s *service) TagCurrentVersion() error {
	currentVersion, err := s.GetCurrentVersion()
	if err != nil {
		return err
	}

	g := git.NewGitService(s.pathToGit)
	err = g.CreateTag(currentVersion)
	if err != nil {
		return err
	}

	err = g.PushTag(currentVersion)
	return err
}
