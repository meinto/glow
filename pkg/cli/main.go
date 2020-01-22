package main

//go:generate ../../.circleci/generate-assets.sh

import (
	"log"

	"github.com/meinto/glow/cmd"
	"github.com/meinto/glow/git"
	clicmd "github.com/meinto/glow/pkg/cli/cmd"
	"github.com/spf13/viper"
)

func main() {
	exec := cmd.NewCmdExecutor("/bin/bash")
	g := git.NewNativeService(exec)
	rootRepoPath, _, _, err := g.GitRepoPath()
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

	clicmd.Execute()
}
