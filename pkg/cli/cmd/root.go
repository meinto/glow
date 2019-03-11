package cmd

import (
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmdOptions struct {
	Author                string
	GitPath               string
	UseBuiltInGitBindings bool
}

var logger log.Logger

var rootCmd = &cobra.Command{
	Use:   "glow",
	Short: "small tool to adapt git-flow for gitlab",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&rootCmdOptions.Author, "author", "a", "name of the author")
	rootCmd.PersistentFlags().StringVar(&rootCmdOptions.GitPath, "gitPath", "/usr/local/bin/git", "path to native git installation")
	rootCmd.PersistentFlags().BoolVar(&rootCmdOptions.UseBuiltInGitBindings, "useBuiltInGitBindings", false, "defines wether build or native in git client should be used.")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("gitPath", rootCmd.PersistentFlags().Lookup("gitPath"))
	viper.BindPFlag("useBuiltInGitBindings", rootCmd.PersistentFlags().Lookup("useBuiltInGitBindings"))
}

func Execute() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
