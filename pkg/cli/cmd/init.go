package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	l "github.com/meinto/glow/logging"
	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
	"github.com/meinto/glow/pkg/cli/cmd/internal/util"
	"github.com/meinto/promter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Author            string `json:"author,omitempty"`
	GitProviderDomain string `json:"gitProviderDomain,omitempty"`
	GitProvider       string `json:"gitProvider,omitempty"`
	ProjectNamespace  string `json:"projectNamespace,omitempty"`
	ProjectName       string `json:"projectName,omitempty"`
}

type PrivateConfig struct {
	Token string `json:"token,omitempty"`
}

var publicConfigFileName = "glow.config.json"
var privateConfigFileName = "glow.private.json"

var initConfig *viper.Viper
var initPrivateConfig *viper.Viper

type InitCommand struct {
	command.Service
}

func (cmd *InitCommand) PostSetup(parent command.Service) command.Service {
	parent.Add(cmd)
	return cmd
}

var initCmd = SetupInitCommand(RootCmd)

func SetupInitCommand(parent command.Service) command.Service {
	return command.Setup(&InitCommand{
		Service: &command.Command{
			Command: &cobra.Command{
				Use:   "init",
				Short: "init glow",
			},
			PreRun: func(cmd command.Service, args []string) {
				rootRepoPath, _, _, err := cmd.GitClient().GitRepoPath()
				if err != nil {
					rootRepoPath = "."
				}

				initConfig = viper.New()
				initConfig.SetConfigName("glow.config")
				initConfig.AddConfigPath(rootRepoPath)
				initConfig.ReadInConfig()

				initPrivateConfig = viper.New()
				initPrivateConfig.SetConfigName("glow.private")
				initPrivateConfig.AddConfigPath(rootRepoPath)
				initPrivateConfig.ReadInConfig()
			},
			Run: func(cmd command.Service, args []string) {
				p := promter.NewPromter()

				author, err := p.TextDefault(
					"Short author name; Will be used for the 'author part' in feature branch names",
					initConfig.GetString("author"),
				)
				util.ExitOnError(err)

				defaultUrl := initConfig.GetString("gitProviderDomain")
				if strings.TrimSpace(defaultUrl) == "" {
					defaultUrl = "https://gitlab.com"
				}
				gitProviderDomain, err := p.URLDefault("Your git host api endpoint", defaultUrl)
				util.ExitOnError(err)

				_, gitProvider, err := p.SelectDefault(
					"Select which git provider you use",
					initConfig.GetString("gitProvider"),
					[]string{"gitlab", "github"},
				)
				util.ExitOnError(err)

				projectNamespace, err := p.TextDefault(
					"Project namespace",
					initConfig.GetString("projectNamespace"),
				)
				util.ExitOnError(err)

				projectName, err := p.TextDefault(
					"Project name",
					initConfig.GetString("projectName"),
				)
				util.ExitOnError(err)

				if _, err := os.Stat(publicConfigFileName); !os.IsNotExist(err) {
					replace, err := replaceFile(publicConfigFileName)
					if err != nil {
						log.Fatal(err)
					}
					if !replace {
						log.Fatal("file not replaced")
					}
				}

				var config = Config{
					Author:            author,
					GitProviderDomain: gitProviderDomain,
					GitProvider:       gitProvider,
					ProjectNamespace:  projectNamespace,
					ProjectName:       projectName,
				}
				writeJSONFile(config, publicConfigFileName)

				token, err := p.Text("Git provider ci token")
				if err != nil {
					log.Fatalf("error setting gitlab ci token: %s", err)
				}
				addToGitIgnore(privateConfigFileName)

				privateConfig := PrivateConfig{
					Token: token,
				}
				writeJSONFile(privateConfig, privateConfigFileName)

				l.Log().Info(l.Fields{
					"config":            config,
					"author":            author,
					"gitProviderDomain": gitProviderDomain,
					"gitProvider":       gitProvider,
					"projectNamespace":  projectNamespace,
					"projectName":       projectName,
					"token":             token,
				})
			},
		},
	}, parent)
}

func writeJSONFile(jsonContent interface{}, fileName string) error {
	newJSONContent, _ := json.MarshalIndent(jsonContent, "", "  ")
	err := ioutil.WriteFile(fileName, newJSONContent, 0644)
	if err != nil {
		return fmt.Errorf("error writing %s: %s", fileName, err.Error())
	}
	return nil
}

func addToGitIgnore(configFileName string) {
	f, err := os.OpenFile(".gitignore", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("\n%s", configFileName)); err != nil {
		panic(err)
	}
}
