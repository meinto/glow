package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	l "github.com/meinto/glow/logging"
	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
	"github.com/meinto/promter"
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
				p := promter.NewPromter()

				author, err := p.Text("Short author name; Will be used for the 'author part' in feature branch names")
				if err != nil {
					log.Fatalf("error setting author: %s", err)
				}

				gitProviderDomain, err := p.URLDefault("Your git host api endpoint (%s)", "https://gitlab.com")
				if err != nil {
					log.Fatalf("error setting git provider api endpoint: %s", err)
				}

				_, gitProvider, err := p.Select(
					"Select which git provider you use",
					[]string{"gitlab", "github"},
				)
				if err != nil {
					log.Fatalf("error setting git provider: %s", err)
				}

				projectNamespace, err := p.Text("Project namespace")
				if err != nil {
					log.Fatalf("error setting project namespace: %s", err)
				}

				projectName, err := p.Text("Project name")
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

				token, err := p.Text("Git provider ci token")
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
	}, parent)
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
