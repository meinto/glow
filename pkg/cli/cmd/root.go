package cmd

import (
	"fmt"
	"log"
	"os"

	"net/url"

	kitlog "github.com/go-kit/kit/log"
	"github.com/gobuffalo/packr/v2"
	"github.com/meinto/glow/pkg/cli/cmd/util"
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

var logger kitlog.Logger

var rootCmd = &cobra.Command{
	Use:     "glow",
	Short:   "small tool to adapt git-flow for gitlab",
	Version: "0.0.0", // needed to set the version dynamically
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if rootCmdOptions.CICDOrigin != "" {
			g, err := util.GetGitClient()
			util.ExitOnError(err)

			util.ExitOnError(g.SetCICDOrigin(rootCmdOptions.CICDOrigin))
		} else if rootCmdOptions.DetectCICDOrigin {
			g, err := util.GetGitClient()
			util.ExitOnError(err)

			cicdOrigin := detectCICDOrigin()
			util.ExitOnError(g.SetCICDOrigin(cicdOrigin))
		}
	},
}

func detectCICDOrigin() string {
	gitProviderDomain := viper.GetString("gitProviderDomain")
	util.CheckRequiredStringField(gitProviderDomain, "gitProviderDomain")

	projectNamespace := viper.GetString("projectNamespace")
	util.CheckRequiredStringField(projectNamespace, "projectNamespace")

	projectName := viper.GetString("projectName")
	util.CheckRequiredStringField(projectName, "projectName")

	gitUser := os.Getenv("CI_GIT_USER")
	gitToken := os.Getenv("CI_GIT_TOKEN")

	if gitUser != "" && gitToken != "" {
		gitProviderURL, err := url.Parse(gitProviderDomain)
		util.ExitOnErrorWithMessage("couldn't parse gitProviderDomain")(err)

		originURL := fmt.Sprintf(
			"%s/%s/%s.git",
			gitProviderURL.Host, projectNamespace, projectName,
		)

		return fmt.Sprintf(
			"%s://%s:%s@%s",
			gitProviderURL.Scheme, gitUser, gitToken, originURL,
		)
	}

	return fmt.Sprintf(
		"%s/%s/%s.git",
		gitProviderDomain, projectNamespace, projectName,
	)
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
	rootCmd.PersistentFlags().StringVar(&rootCmdOptions.CICDOrigin, "cicdOrigin", "", "provide a git origin url where a pipeline can push things via token")
	rootCmd.PersistentFlags().BoolVar(&rootCmdOptions.DetectCICDOrigin, "detectCicdOrigin", false, "auto detect a git origin url where a pipeline can push things via token")
	rootCmd.PersistentFlags().BoolVar(&rootCmdOptions.CI, "ci", false, "detects if command is running in a ci")
	rootCmd.PersistentFlags().BoolVar(&rootCmdOptions.SkipChecks, "skipChecks", false, "skip checks like accidentally creating git-flow branches from wrong source branch")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("gitPath", rootCmd.PersistentFlags().Lookup("gitPath"))
}

func Execute() {
	logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
