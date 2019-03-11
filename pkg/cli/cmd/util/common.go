package util

import (
	"log"
	"os"

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
