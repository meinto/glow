package cmd

import (
	"log"
	"os"

	"github.com/gobuffalo/packr/v2"
	l "github.com/meinto/glow/logging"
	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
	"github.com/meinto/glow/pkg/cli/cmd/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmdOptions struct {
	Author           string
	GitPath          string
	CICDOrigin       string
	DetectCICDOrigin bool
	CI               bool
	SkipChecks       bool
}

var RootCmd = &command.Command{
	Command: &cobra.Command{
		Use:     "glow",
		Short:   "small tool to adapt git-flow for gitlab",
		Version: "0.0.0", // needed to set the version dynamically
	},
	PersistentPreRun: func(cmd command.Service, args []string) {
		if rootCmdOptions.CICDOrigin != "" {
			util.ExitOnError(cmd.GitClient().SetCICDOrigin(rootCmdOptions.CICDOrigin))
		} else if rootCmdOptions.DetectCICDOrigin {
			cicdOrigin, err := cmd.GitProvider().DetectCICDOrigin()
			util.ExitOnError(err)
			util.ExitOnError(cmd.GitClient().SetCICDOrigin(cicdOrigin))
		}
	},
	Run: func(cmd command.Service, args []string) {
		log.Println("hello from glow")
	},
}

func init() {
	box := packr.New("build-assets", "../../../buildAssets")
	version, err := box.FindString("VERSION")
	if err != nil {
		l.Log().Error(err)
		version = "0.0.0"
	}
	RootCmd.SetVersionTemplate(version)

	RootCmd.PersistentFlags().StringVarP(&rootCmdOptions.Author, "author", "a", "", "name of the author")
	RootCmd.PersistentFlags().StringVar(&rootCmdOptions.GitPath, "gitPath", "/usr/local/bin/git", "path to native git installation")
	RootCmd.PersistentFlags().StringVar(&rootCmdOptions.CICDOrigin, "cicdOrigin", "", "provide a git origin url where a pipeline can push things via token")
	RootCmd.PersistentFlags().BoolVar(&rootCmdOptions.DetectCICDOrigin, "detectCicdOrigin", false, "auto detect a git origin url where a pipeline can push things via token")
	RootCmd.PersistentFlags().BoolVar(&rootCmdOptions.CI, "ci", false, "detects if command is running in a ci")
	RootCmd.PersistentFlags().BoolVar(&rootCmdOptions.SkipChecks, "skipChecks", false, "skip checks like accidentally creating git-flow branches from wrong source branch")
	viper.BindPFlag("author", RootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("gitPath", RootCmd.PersistentFlags().Lookup("gitPath"))
}

func Execute() {
	if err := RootCmd.
		Init().
		Patch().
		Execute(); err != nil {
		l.Log().Error(err)
		os.Exit(1)
	}
}
