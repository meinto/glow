package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

func CheckRequiredStringField(val, fieldName string) {
	if val == "" {
		log.Fatalf("please provide %s", fieldName)
	}
}

func CreateMergeRequest(rawSource, rawTarget string) {
	source := strings.TrimPrefix(rawSource, "refs/heads/")
	target := strings.TrimPrefix(rawTarget, "refs/heads/")
	CheckRequiredStringField(source, "source branch")
	CheckRequiredStringField(target, "target branch")
	CheckRequiredStringField(viper.GetString("gitlabEndpoint"), "gitlab endpoint")
	CheckRequiredStringField(viper.GetString("projectNamespace"), "project namespace")
	CheckRequiredStringField(viper.GetString("projectName"), "project name")
	CheckRequiredStringField(viper.GetString("gitlabCIToken"), "gitlab ci token")

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
}
