package util

import (
	"os"

	"github.com/spf13/viper"
)

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
