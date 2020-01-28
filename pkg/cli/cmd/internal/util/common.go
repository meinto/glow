package util

import (
	"errors"
	"log"

	"github.com/meinto/glow/cmd"

	"github.com/meinto/glow/gitprovider"

	"github.com/meinto/glow/git"
	"github.com/spf13/viper"
)

func CheckRequiredStringField(val, fieldName string) {
	if val == "" {
		log.Fatalf("please provide %s", fieldName)
	}
}

func GetGitClient() (git.Service, error) {
	var s git.Service
	exec := cmd.NewCmdExecutor("/bin/bash")
	s = git.NewNativeService(exec)
	s = git.NewLoggingService(s)
	return s, nil
}

func GetGitProvider() (gitprovider.Service, error) {
	var s gitprovider.Service
	switch viper.GetString("gitProvider") {
	case "github":
		s = gitprovider.NewGithubService(
			viper.GetString("gitProviderDomain"),
			viper.GetString("projectNamespace"),
			viper.GetString("projectName"),
			viper.GetString("token"),
		)
		break
	case "gitlab":
		s = gitprovider.NewGitlabService(
			viper.GetString("gitProviderDomain"),
			viper.GetString("projectNamespace"),
			viper.GetString("projectName"),
			viper.GetString("token"),
		)
		break
	default:
		return nil, errors.New("git provider not specified")
	}
	return s, nil
}
