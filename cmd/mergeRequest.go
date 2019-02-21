package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var mergeRequestCmdOptions struct {
	GitlabEndpoint   string
	ProjectNamespace string
	ProjectName      string
	GitlabCIToken    string
}

func init() {
	rootCmd.AddCommand(mergeRequestCmd)
	mergeRequestCmd.Flags().StringVarP(&mergeRequestCmdOptions.GitlabEndpoint, "endpoint", "e", "", "gitlab endpoint")
	mergeRequestCmd.Flags().StringVarP(&mergeRequestCmdOptions.ProjectNamespace, "namespace", "n", "", "project namespace")
	mergeRequestCmd.Flags().StringVarP(&mergeRequestCmdOptions.ProjectName, "project", "p", "", "project name")
	mergeRequestCmd.Flags().StringVarP(&mergeRequestCmdOptions.GitlabCIToken, "token", "t", "", "gitlab ci token")
	viper.BindPFlag("gitlabEndpoint", mergeRequestCmd.Flags().Lookup("endpoint"))
	viper.BindPFlag("projectNamespace", mergeRequestCmd.Flags().Lookup("namespace"))
	viper.BindPFlag("projectName", mergeRequestCmd.Flags().Lookup("project"))
	viper.BindPFlag("gitlabCIToken", mergeRequestCmd.Flags().Lookup("token"))
}

var mergeRequestCmd = &cobra.Command{
	Use:   "mergeRequest",
	Short: "create a merge request on gitlab",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Fatal("Please provide source and target branch")
		}

		source := args[0]
		target := args[1]

		CheckRequiredStringField(source, "source branch")
		CheckRequiredStringField(target, "target branch")
		CheckRequiredStringField(mergeRequestCmdOptions.GitlabEndpoint, "gitlab endpoint")
		CheckRequiredStringField(mergeRequestCmdOptions.ProjectNamespace, "project namespace")
		CheckRequiredStringField(mergeRequestCmdOptions.ProjectName, "project name")
		CheckRequiredStringField(mergeRequestCmdOptions.GitlabCIToken, "gitlab ci token")

		type Payload struct {
			SourceBranch       string `json:"source_branch"`
			TargetBranch       string `json:"target_branch"`
			Title              string `json:"title"`
			RemoveSourceBranch bool   `json:"remove_source_branch"`
		}

		data := Payload{
			SourceBranch:       source,
			TargetBranch:       target,
			Title:              fmt.Sprintf("Merge %s in %s", source, target),
			RemoveSourceBranch: false,
		}
		payloadBytes, err := json.Marshal(data)
		CheckForError(err, "mergeRequestCmd marshal json")

		body := bytes.NewReader(payloadBytes)

		requestURI := fmt.Sprintf(
			"%s/api/v4/projects/%s%s%s/merge_requests",
			viper.GetString("gitlabEndpoint"),
			viper.GetString("projectNamespace"),
			"%2F",
			viper.GetString("projectName"),
		)
		req, err := http.NewRequest("POST", requestURI, body)
		CheckForError(err, "mergeRequestCmd prepare request")

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Private-Token", viper.GetString("gitlabCIToken"))

		resp, err := http.DefaultClient.Do(req)
		CheckForError(err, "mergeRequestCmd do request")

		defer resp.Body.Close()

	},
}
