package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmdOptions struct {
	Author               string
	GitPath              string
	UseNativeGitBindings []string
}

var rootCmd = &cobra.Command{
	Use:   "glow",
	Short: "small tool to adapt git-flow for gitlab",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootCmdOptions.Author, "author", "a", "test", "name of the author")
	rootCmd.PersistentFlags().StringVar(&rootCmdOptions.GitPath, "gitPath", "/usr/local/bin/git", "path to native git installation")
	rootCmd.PersistentFlags().StringArrayVar(&rootCmdOptions.UseNativeGitBindings, "useNativeGitBindings", []string{}, "defines which git actions should be performed with the native git client")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("gitPath", rootCmd.PersistentFlags().Lookup("gitPath"))
	viper.BindPFlag("useNativeGitBindings", rootCmd.PersistentFlags().Lookup("useNativeGitBindings"))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
