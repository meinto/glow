package main

//go:generate ../../.circleci/generate-assets.sh

import (
	"github.com/meinto/glow/cmd"
	"github.com/meinto/glow/git"
	l "github.com/meinto/glow/logging"
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
	l.Log().Warn(l.Fields{"msg": "there is no glow config"})

	viper.SetConfigName("glow.private")
	viper.AddConfigPath(rootRepoPath)
	err = viper.MergeInConfig()
	l.Log().Warn(l.Fields{"msg": "there is no private glow config"})

	viper.SetEnvPrefix("glow")
	err = viper.BindEnv("token")
	l.Log().Warn(l.Fields{"msg": "env GLOW_TOKEN is missing"})

	clicmd.Execute()
}
