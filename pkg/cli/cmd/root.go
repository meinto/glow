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

func Registry(rootCmd command.Service) {
	ReleaseCmd(rootCmd)
}

func SetupRootCmd(cmd command.Service) command.Service {
	box := packr.New("build-assets", "../../../buildAssets")
	version, err := box.FindString("VERSION")
	if err != nil {
		l.Log().Error(err)
		version = "0.0.0"
	}
	cmd.Cmd().SetVersionTemplate(version)

	cmd.Cmd().PersistentFlags().StringVarP(&rootCmdOptions.Author, "author", "a", "", "name of the author")
	cmd.Cmd().PersistentFlags().StringVar(&rootCmdOptions.GitPath, "gitPath", "/usr/local/bin/git", "path to native git installation")
	cmd.Cmd().PersistentFlags().StringVar(&rootCmdOptions.CICDOrigin, "cicdOrigin", "", "provide a git origin url where a pipeline can push things via token")
	cmd.Cmd().PersistentFlags().BoolVar(&rootCmdOptions.DetectCICDOrigin, "detectCicdOrigin", false, "auto detect a git origin url where a pipeline can push things via token")
	cmd.Cmd().PersistentFlags().BoolVar(&rootCmdOptions.CI, "ci", false, "detects if command is running in a ci")
	cmd.Cmd().PersistentFlags().BoolVar(&rootCmdOptions.SkipChecks, "skipChecks", false, "skip checks like accidentally creating git-flow branches from wrong source branch")
	viper.BindPFlag("author", cmd.Cmd().PersistentFlags().Lookup("author"))
	viper.BindPFlag("gitPath", cmd.Cmd().PersistentFlags().Lookup("gitPath"))

	return cmd
}

func CreateRootCmd() command.Service {
	return SetupRootCmd(&command.Command{
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
	})
}

var RootCmd = CreateRootCmd()

func init() {
	Registry(RootCmd)
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
