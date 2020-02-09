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

var RootCmdOptions struct {
	Author           string
	GitPath          string
	CICDOrigin       string
	DetectCICDOrigin bool
	CI               bool
	SkipChecks       bool
}

type RootCommand struct {
	command.Service
}

func (cmd *RootCommand) PostSetup(parent command.Service) command.Service {
	box := packr.New("build-assets", "../../../buildAssets")
	version, err := box.FindString("VERSION")
	if err != nil {
		l.Log().Error(err)
		version = "0.0.0"
	}
	cmd.Cmd().SetVersionTemplate(version)

	cmd.Cmd().PersistentFlags().StringVarP(&RootCmdOptions.Author, "author", "a", "", "name of the author")
	cmd.Cmd().PersistentFlags().StringVar(&RootCmdOptions.GitPath, "gitPath", "/usr/local/bin/git", "path to native git installation")
	cmd.Cmd().PersistentFlags().StringVar(&RootCmdOptions.CICDOrigin, "cicdOrigin", "", "provide a git origin url where a pipeline can push things via token")
	cmd.Cmd().PersistentFlags().BoolVar(&RootCmdOptions.DetectCICDOrigin, "detectCicdOrigin", false, "auto detect a git origin url where a pipeline can push things via token")
	cmd.Cmd().PersistentFlags().BoolVar(&RootCmdOptions.CI, "ci", false, "detects if command is running in a ci")
	cmd.Cmd().PersistentFlags().BoolVar(&RootCmdOptions.SkipChecks, "skipChecks", false, "skip checks like accidentally creating git-flow branches from wrong source branch")
	viper.BindPFlag("author", cmd.Cmd().PersistentFlags().Lookup("author"))
	viper.BindPFlag("gitPath", cmd.Cmd().PersistentFlags().Lookup("gitPath"))
	return cmd
}

var RootCmd = SetupRootCommand()

func SetupRootCommand() command.Service {
	return command.Setup(&RootCommand{
		&command.Command{
			Command: &cobra.Command{
				Use:     "glow",
				Short:   "small tool to adapt git-flow for gitlab",
				Version: "0.0.0", // needed to set the version dynamically
			},
			PersistentPreRun: func(cmd command.Service, args []string) {
				if RootCmdOptions.CICDOrigin != "" {
					util.ExitOnError(cmd.GitClient().SetCICDOrigin(RootCmdOptions.CICDOrigin))
				} else if RootCmdOptions.DetectCICDOrigin {
					cicdOrigin, err := cmd.GitProvider().DetectCICDOrigin()
					util.ExitOnError(err)
					util.ExitOnError(cmd.GitClient().SetCICDOrigin(cicdOrigin))
				}
			},
			Run: func(cmd command.Service, args []string) {
				log.Println("hello from glow")
			},
		},
	}, nil)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		l.Log().Error(err)
		os.Exit(1)
	}
}
