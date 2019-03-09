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

	viper.SetConfigName("glow")
	viper.AddConfigPath(rootRepoPath)
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		log.Println("there is no glow config")
	}

	cmd.Execute()
}
