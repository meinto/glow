package cmd

import (
	"fmt"
	"log"
	"os"

	kitlog "github.com/go-kit/kit/log"
	"github.com/gobuffalo/packr/v2"
	"github.com/meinto/glow/pkg/cli/cmd/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmdOptions struct {
	Author                string
	GitPath               string
	UseBuiltInGitBindings bool
	CICDOrigin            string
	CI                    bool
	SkipChecks            bool
}

var logger kitlog.Logger

var rootCmd = &cobra.Command{
	Use:     "glow",
	Short:   "small tool to adapt git-flow for gitlab",
	Version: "0.0.0", // needed to set the version dynamically
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if rootCmdOptions.CICDOrigin != "" {
			g, err := util.GetGitClient()
			util.CheckForError(err, "GetGitClient")

			err = g.SetCICDOrigin(rootCmdOptions.CICDOrigin)
			util.CheckForError(err, "SetCICDOrigin")
		}
	},
}

func init() {
	box := packr.New("build-assets", "../../../buildAssets")
	version, err := box.FindString("VERSION")
	if err != nil {
		log.Println(err)
		version = "0.0.0"
	}
	rootCmd.SetVersionTemplate(version)

	rootCmd.PersistentFlags().StringVarP(&rootCmdOptions.Author, "author", "a", "", "name of the author")
	rootCmd.PersistentFlags().StringVar(&rootCmdOptions.GitPath, "gitPath", "/usr/local/bin/git", "path to native git installation")
	rootCmd.PersistentFlags().BoolVar(&rootCmdOptions.UseBuiltInGitBindings, "useBuiltInGitBindings", false, "defines wether build or native in git client should be used.")
	rootCmd.PersistentFlags().StringVar(&rootCmdOptions.CICDOrigin, "cicdOrigin", "", "provide a git origin url where a pipeline can push things via token")
	rootCmd.PersistentFlags().BoolVar(&rootCmdOptions.CI, "ci", false, "detects if command is running in a ci")
	rootCmd.PersistentFlags().BoolVar(&rootCmdOptions.SkipChecks, "skipChecks", false, "skip checks like accidentally creating git-flow branches from wrong source branch")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("gitPath", rootCmd.PersistentFlags().Lookup("gitPath"))
	viper.BindPFlag("useBuiltInGitBindings", rootCmd.PersistentFlags().Lookup("useBuiltInGitBindings"))
}

func Execute() {
	logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
