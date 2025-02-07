package pkg

import (
	"log"
	"os"
	"os/user"

	"github.com/meinto/glow/cmd"
	"github.com/meinto/glow/git"
	l "github.com/meinto/glow/logging"
	"github.com/spf13/viper"
)

var configInitialized = false

func InitGlobalConfig() {
	if !configInitialized {
		configInitialized = true
		exec := cmd.NewCmdExecutor("/bin/bash")
		g := git.NewNativeService(git.Options{
			CmdExecutor: exec,
		})
		rootRepoPath, _, _, err := g.GitRepoPath()
		if err != nil {
			rootRepoPath = "."
		}

		viper.SetDefault("versionFile", "VERSION")
		viper.SetDefault("logLevel", "info")
		viper.SetDefault("versionFileType", "raw")
		viper.SetDefault("mainBranch", "master")
		viper.SetDefault("devBranch", "develop")

		viper.SetConfigName("glow.config")
		viper.AddConfigPath(rootRepoPath)
		err = viper.ReadInConfig()
		l.Configure(l.Options{
			l.GetLevel(viper.GetString("logLevel")),
		})

		l.Log().
			Info(viper.AllSettings()).
			ErrorFields(err, l.Fields{"msg": "there is no glow config"})

		viper.SetConfigName("glow.private")
		viper.AddConfigPath(rootRepoPath)
		err = viper.MergeInConfig()
		l.Log().ErrorFields(err, l.Fields{"msg": "there is no private glow config"})

		// load global config
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		globalConfigPath := usr.HomeDir + "/.glow"
		viper.SetConfigName("glow.config")
		viper.AddConfigPath(globalConfigPath)
		err = viper.MergeInConfig()
		l.Log().ErrorFields(err, l.Fields{"msg": "there is no global glow config"})

		viper.SetConfigName("glow.private")
		viper.AddConfigPath(globalConfigPath)
		err = viper.MergeInConfig()
		l.Log().ErrorFields(err, l.Fields{"msg": "there is no global private glow config"})

		viper.SetEnvPrefix("glow")
		err = viper.BindEnv("token")
		l.Log().
			WarnIf(l.Fields{"msg": "env GLOW_TOKEN is missing"}, os.Getenv("GLOW_TOKEN") == "").
			Error(err)
	}
}
