package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

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
	SquashCommits     bool   `json:"squashCommits,omitempty"`
	VersionFile       string `json:"versionFile,omitempty"`
	VersionFileType   string `json:"versionFileType,omitempty"`
	LogLevel          string `json:"logLevel,omitempty"`
}

type PrivateConfig struct {
	Token string `json:"token,omitempty"`
}

var configDir = "."
var initGlobal = false
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
				p := promter.NewPromter()
				index := initProjectOrGlobalConfig(p)
				if index == 1 {
					usr, err := user.Current()
					if err != nil {
						log.Fatal(err)
					}
					configDir = usr.HomeDir + "/.glow"
					initGlobal = true
				} else {
					rootRepoPath, _, _, err := cmd.GitClient().GitRepoPath()
					if err != nil {
						rootRepoPath = "."
					}
					configDir = rootRepoPath
					initGlobal = false
				}

				initConfig = viper.New()
				initConfig.SetConfigName("glow.config")
				initConfig.AddConfigPath(configDir)
				initConfig.ReadInConfig()

				initPrivateConfig = viper.New()
				initPrivateConfig.SetConfigName("glow.private")
				initPrivateConfig.AddConfigPath(configDir)
				initPrivateConfig.ReadInConfig()
			},
			Run: func(cmd command.Service, args []string) {
				p := promter.NewPromter()

				var config = Config{}
				author(p, &config)
				gitProviderDomain(p, &config)
				gitProvider(p, &config)
				if !initGlobal {
					projectNamespace(p, &config)
					projectName(p, &config)
					squashCommits(p, &config)
					versionFile(p, &config)
					versionFileType(p, &config)
				}
				logLevel(p, &config)

				promtReplaceFile(configDir + "/" + publicConfigFileName)

				writeJSONFile(config, configDir+"/"+publicConfigFileName)

				shouldCreate := shouldCreatePrivateConfig(p)

				if shouldCreate {
					token, err := p.Text("Git provider ci token")
					if err != nil {
						log.Fatalf("error setting gitlab ci token: %s", err)
					}
					addToGitIgnore(privateConfigFileName)

					privateConfig := PrivateConfig{
						Token: token,
					}

					promtReplaceFile(configDir + "/" + privateConfigFileName)

					writeJSONFile(privateConfig, configDir+"/"+privateConfigFileName)
				}
			},
		},
	}, parent)
}

func writeJSONFile(jsonContent interface{}, fileName string) error {
	newJSONContent, _ := json.MarshalIndent(jsonContent, "", "  ")
	os.MkdirAll(filepath.Dir(fileName), os.ModePerm)
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

func shouldCreatePrivateConfig(p promter.Promter) bool {
	index, _, err := p.YesNo("Do you want to create a private.glow config?")
	if err != nil {
		return false
	}
	if index == 0 {
		return true
	}
	return false
}

func initProjectOrGlobalConfig(p promter.Promter) int {
	index, _, _ := p.Select("Do you want a global or a project setup?", []string{
		"project", "global",
	})
	return index
}

// promts
func author(p promter.Promter, config *Config) {
	val, err := p.TextDefault(
		"Short author name; Will be used for the 'author part' in feature branch names",
		initConfig.GetString("author"),
	)
	util.ExitOnError(err)
	config.Author = val
}

func gitProviderDomain(p promter.Promter, config *Config) {
	defaultUrl := initConfig.GetString("gitProviderDomain")
	if strings.TrimSpace(defaultUrl) == "" {
		defaultUrl = "https://gitlab.com"
	}
	val, err := p.URLDefault("Your git host api endpoint", defaultUrl)
	util.ExitOnError(err)
	config.GitProviderDomain = val
}

func gitProvider(p promter.Promter, config *Config) {
	_, val, err := p.SelectDefault(
		"Select which git provider you use",
		initConfig.GetString("gitProvider"),
		[]string{"gitlab", "github"},
	)
	util.ExitOnError(err)
	config.GitProvider = val
}

func projectNamespace(p promter.Promter, config *Config) {
	val, err := p.TextDefault(
		"Project namespace",
		initConfig.GetString("projectNamespace"),
	)
	util.ExitOnError(err)
	config.ProjectNamespace = val
}

func projectName(p promter.Promter, config *Config) {
	val, err := p.TextDefault(
		"Project name",
		initConfig.GetString("projectName"),
	)
	util.ExitOnError(err)
	config.ProjectName = val
}

func squashCommits(p promter.Promter, config *Config) {
	defaultVal := "Yes"
	if initConfig.GetBool("squashCommits") == false {
		defaultVal = "No"
	}
	index, _, err := p.YesNoDefault("Squash commits?", defaultVal)
	util.ExitOnError(err)
	switch index {
	case 0:
		config.SquashCommits = true
	case 1:
		config.SquashCommits = false
	default:
		config.SquashCommits = false
	}
}

func versionFile(p promter.Promter, config *Config) {
	val, err := p.TextDefault(
		"Version file name",
		initConfig.GetString("versionFile"),
	)
	util.ExitOnError(err)
	config.VersionFile = val
}

func versionFileType(p promter.Promter, config *Config) {
	val, err := p.TextDefault(
		"Version file type",
		initConfig.GetString("versionFileType"),
	)
	util.ExitOnError(err)
	config.VersionFileType = val
}

func logLevel(p promter.Promter, config *Config) {
	_, val, err := p.SelectDefault(
		"Log level",
		initConfig.GetString("logLevel"),
		[]string{
			"trace",
			"debug",
			"panic",
			"fatal",
			"error",
			"warning",
			"info",
		},
	)
	util.ExitOnError(err)
	config.LogLevel = val
}
