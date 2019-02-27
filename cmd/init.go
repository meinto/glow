package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type ConfigType struct {
	Author           string `json:"author,omitempty"`
	GitlabEndpoint   string `json:"gitlab_endpoint,omitempty"`
	ProjectNamespace string `json:"project_namespace,omitempty"`
	ProjectName      string `json:"project_name,omitempty"`
	GitlabCIToken    string `json:"gitlab_ci_token,omitempty"`
}

var configFileName = "glow.json"

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init glow",
	Run: func(cmd *cobra.Command, args []string) {

		author, err := promtNotEmpty("Your name - Author ")
		if err != nil {
			log.Fatalf("error setting author: %s", err)
		}

		gitlabEndpoint, err := promptURL("Your gitlab endpoint (%s) ", "https://gitlab.com")
		if err != nil {
			log.Fatalf("error setting gitlab endpoint: %s", err)
		}

		projectNamespace, err := promtNotEmpty("Project namespace ")
		if err != nil {
			log.Fatalf("error setting project namespace: %s", err)
		}

		projectName, err := promtNotEmpty("Project name ")
		if err != nil {
			log.Fatalf("error setting project name: %s", err)
		}

		gitlabCIToken, err := promtNotEmpty("Gitlab ci token ")
		if err != nil {
			log.Fatalf("error setting gitlab ci token: %s", err)
		}

		if _, err := os.Stat(configFileName); !os.IsNotExist(err) {
			replace, err := replaceFile(configFileName)
			if err != nil {
				log.Fatal(err)
			}
			if !replace {
				log.Fatal("file not replaced")
			}
		}

		addToGitIgnore(configFileName)

		var config = ConfigType{
			Author:           author,
			GitlabEndpoint:   gitlabEndpoint,
			ProjectNamespace: projectNamespace,
			ProjectName:      projectName,
			GitlabCIToken:    gitlabCIToken,
		}
		writeJSONFile(config, configFileName)

		log.Println(config, author, gitlabEndpoint, projectNamespace, projectName, gitlabCIToken)
	},
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
