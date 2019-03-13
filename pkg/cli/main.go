package main

import (
	"log"
	"os"

	kitlog "github.com/go-kit/kit/log"

	"github.com/meinto/glow/git"
	"github.com/meinto/glow/pkg/cli/cmd"
	"github.com/spf13/viper"
)

func main() {
	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)

	g := git.NewGoGitService()
	g = git.NewLoggingService(logger, g)
	rootRepoPath, err := g.GitRepoPath()
	if err != nil {
		rootRepoPath = "."
	}

	viper.SetConfigName("glow.config")
	viper.AddConfigPath(rootRepoPath)
	err = viper.ReadInConfig()
	if err != nil {
		log.Println("there is no glow config")
	}

	viper.SetConfigName("glow.private")
	viper.AddConfigPath(rootRepoPath)
	err = viper.MergeInConfig()
	if err != nil {
		log.Println("there is no private glow config")
	}

	viper.SetEnvPrefix("glow")
	err = viper.BindEnv("token")
	if err != nil {
		log.Println("env GLOW_TOKEN is missing")
	}

	cmd.Execute()
}
