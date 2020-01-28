package util

import (
	"os"

	"github.com/spf13/viper"

	"github.com/meinto/glow/semver"
)

func ProcessVersion(versionArg, versionFile, versionFileType string) (string, semver.Service) {
	version := versionArg

	g, err := GetGitClient()
	ExitOnError(err)

	pathToRepo, _, _, err := g.GitRepoPath()
	ExitOnError(err)

	s := semver.NewSemverService(
		pathToRepo,
		"/bin/bash",
		versionFile,
		versionFileType,
	)

	if version == "current" {
		v, err := s.GetCurrentVersion()
		ExitOnError(err)
		version = v
	}

	if IsSemanticVersion(version) {
		v, err := s.GetNextVersion(version)
		ExitOnError(err)
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
