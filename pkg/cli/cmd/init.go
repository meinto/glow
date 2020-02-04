package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	cobraUtils "github.com/meinto/cobra-utils"
	l "github.com/meinto/glow/logging"
	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
	"github.com/spf13/cobra"
)

type ConfigType struct {
	Author            string `json:"author,omitempty"`
	GitProviderDomain string `json:"gitProviderDomain,omitempty"`
	GitProvider       string `json:"gitProvider,omitempty"`
	ProjectNamespace  string `json:"projectNamespace,omitempty"`
	ProjectName       string `json:"projectName,omitempty"`
	Token             string `json:"token,omitempty"`
}

var publicConfigFileName = "glow.config.json"
var privateConfigFileName = "glow.private.json"

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
		&command.Command{
			Command: &cobra.Command{
				Use:   "init",
				Short: "init glow",
			},
			Run: func(cmd command.Service, args []string) {

				author, err := promtNotEmpty("Short author name; Will be used for the 'author part' in feature branch names")
				if err != nil {
					log.Fatalf("error setting author: %s", err)
				}

				gitProviderDomain, err := promptURL("Your git host api endpoint (%s)", "https://gitlab.com")
				if err != nil {
					log.Fatalf("error setting git provider api endpoint: %s", err)
				}

				_, gitProvider, err := cobraUtils.PromptSelect(
					"Select which git provider you use",
					[]string{"gitlab", "github"},
				)
				if err != nil {
					log.Fatalf("error setting git provider: %s", err)
				}

				projectNamespace, err := promtNotEmpty("Project namespace")
				if err != nil {
					log.Fatalf("error setting project namespace: %s", err)
				}

				projectName, err := promtNotEmpty("Project name")
				if err != nil {
					log.Fatalf("error setting project name: %s", err)
				}

				if _, err := os.Stat(publicConfigFileName); !os.IsNotExist(err) {
					replace, err := replaceFile(publicConfigFileName)
					if err != nil {
						log.Fatal(err)
					}
					if !replace {
						log.Fatal("file not replaced")
					}
				}

				var config = ConfigType{
					Author:            author,
					GitProviderDomain: gitProviderDomain,
					GitProvider:       gitProvider,
					ProjectNamespace:  projectNamespace,
					ProjectName:       projectName,
				}
				writeJSONFile(config, publicConfigFileName)

				token, err := promtNotEmpty("Git provider ci token")
				if err != nil {
					log.Fatalf("error setting gitlab ci token: %s", err)
				}
				addToGitIgnore(privateConfigFileName)

				config = ConfigType{
					Token: token,
				}
				writeJSONFile(config, privateConfigFileName)

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
	})
}

func promtNotEmpty(label string) (string, error) {
	validate := func(input string) error {
		if input == "" {
			return fmt.Errorf("please enter your name")
		}
		return nil
	}

	getNotEmptyValue := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}

	notEmpty, err := getNotEmptyValue.Run()
	if err != nil {
		return "", err
	}

	return strings.ToLower(notEmpty), nil
}

func promptURL(label, defaultValue string) (string, error) {
	validate := func(input string) error {
		if input != "" &&
			!strings.HasPrefix(input, "http://") &&
			!strings.HasPrefix(input, "https://") {
			return fmt.Errorf("please enter a valid url")
		}
		return nil
	}

	getUrl := promptui.Prompt{
		Label:    fmt.Sprintf(label, defaultValue),
		Validate: validate,
	}

	url, err := getUrl.Run()
	if err != nil {
		return "", err
	}

	if url == "" {
		url = defaultValue
	}

	return strings.ToLower(url), nil
}

func writeJSONFile(jsonContent ConfigType, fileName string) error {
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
