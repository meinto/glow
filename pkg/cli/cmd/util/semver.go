package util

import (
	"os"

	"github.com/spf13/viper"

	"github.com/meinto/glow/semver"
)

func ProcessVersion(versionArg, versionFile, versionFileType string) (string, semver.Service) {
	version := versionArg

	g, err := GetGitClient()
	CheckForError(err, "GetGitClient")

	var s semver.Service
	if IsSemanticVersion(version) {
		pathToRepo, err := g.GitRepoPath()
		CheckForError(err, "semver GitRepoPath")
		s = semver.NewSemverService(
			pathToRepo,
			"/bin/bash",
			versionFile,
			versionFileType,
		)
		v, err := s.GetCurrentVersion()
		CheckForError(err, "semver GetNextVersion")
		version = v
	}

	return version, s
}

func HasSemverConfig() bool {
	if _, err := os.Stat(viper.GetString("gitPath") + "/semver.config.json"); os.IsNotExist(err) {
		return false
	}
	return true
}

func IsSemanticVersion(version string) bool {
	if version == "major" || version == "minor" || version == "patch" {
		return true
	}
	return false
}
