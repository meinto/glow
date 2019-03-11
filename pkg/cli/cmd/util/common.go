package util

import (
	"errors"
	"log"
	"os"

	"github.com/meinto/glow/gitprovider"

	kitlog "github.com/go-kit/kit/log"
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
	if viper.GetBool("useBuiltInGitBindings") {
		s = git.NewGoGitService()
	}
	s, err := git.NewNativeService(viper.GetString("gitPath"))
	if err != nil {
		return nil, err
	}
	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
	s = git.NewLoggingService(logger, s)
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
			viper.GetString("gitlabCIToken"),
		)
		break
	default:
		return nil, errors.New("git provider not specified")
	}
	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
	s = gitprovider.NewLoggingService(logger, s)
	return s, nil
}
