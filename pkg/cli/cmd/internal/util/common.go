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
	s = git.NewNativeService(git.Options{
		CmdExecutor: exec,
	})
	return s, nil
}

func GetGitProvider() (gitprovider.Service, error) {
	var s gitprovider.Service
	endpoint := viper.GetString("gitProviderDomain")
	projectNamespace := viper.GetString("projectNamespace")
	projectName := viper.GetString("projectName")
	gitproviderToken := viper.GetString("token")
	switch viper.GetString("gitProvider") {
	case "github":
		s = gitprovider.NewGithubService(gitprovider.Options{
			Endpoint:  endpoint,
			Namespace: projectName,
			Project:   projectNamespace,
			Token:     gitproviderToken,
		})
		break
	case "gitlab":
		s = gitprovider.NewGitlabService(gitprovider.Options{
			Endpoint:  endpoint,
			Namespace: projectName,
			Project:   projectNamespace,
			Token:     gitproviderToken,
		})
		break
	default:
		return nil, errors.New("git provider not specified")
	}
	return s, nil
}
