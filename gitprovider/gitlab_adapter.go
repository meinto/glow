package gitprovider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/meinto/glow"
	"github.com/meinto/glow/git"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type gitlabAdapter struct {
	service
}

func (s *gitlabAdapter) GitService() (gs git.Service) {
	return s.gitService
}

func (s *gitlabAdapter) HTTPClient() HttpClient {
	return s.httpClient
}

func (s *gitlabAdapter) SetHTTPClient(client HttpClient) {
	s.SetHTTPClient(client)
}

func (s *gitlabAdapter) SetGitService(gs git.Service) {
	s.gitService = gs
}

func (a *gitlabAdapter) Close(b glow.Branch) error {
	if b.CanBeClosed() {
		exists, _, _, err := a.GitService().RemoteBranchExists(b.ShortBranchName())
		if exists && err == nil {
			branchList, _, _, err := a.GitService().BranchList()
			if err != nil {
				return errors.Wrap(err, "error getting branch list")
			}
			targets := b.CloseBranches(branchList)

			for _, t := range targets {
				log.Println(t.BranchName())
				err := a.createMergeRequest(b, t, false)
				if err != nil {
					return errors.Wrap(err, "error creating merge request")
				}
			}
			return nil
		}
	}
	return errors.New("cannot be closed")
}

func (a *gitlabAdapter) Publish(b glow.Branch) error {
	if b.CanBePublished() {
		exists, _, _, err := a.GitService().RemoteBranchExists(b.ShortBranchName())
		if exists && err == nil {
			t := b.PublishBranch()
			return a.createMergeRequest(b, t, false)
		}
	}
	return errors.New("cannot be published")
}

func (a *gitlabAdapter) createMergeRequest(source glow.Branch, target glow.Branch, removeSourceBranch bool) error {
	type Payload struct {
		SourceBranch       string `json:"source_branch"`
		TargetBranch       string `json:"target_branch"`
		Title              string `json:"title"`
		RemoveSourceBranch bool   `json:"remove_source_branch"`
		Squash             bool   `json:"squash"`
	}

	type mergeRequestResponse struct {
		URL string `json:"web_url"`
	}

	sourceBranchName := source.ShortBranchName()
	targetBranchName := target.ShortBranchName()

	data := Payload{
		SourceBranch:       sourceBranchName,
		TargetBranch:       targetBranchName,
		Title:              fmt.Sprintf("Merge %s in %s", sourceBranchName, targetBranchName),
		RemoveSourceBranch: removeSourceBranch,
		Squash:             viper.GetBool("squashCommits"),
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "marshal json")
	}

	body := bytes.NewReader(payloadBytes)

	requestURI := fmt.Sprintf(
		"%s/api/v4/projects/%s/merge_requests",
		a.endpoint,
		url.QueryEscape(a.namespace+"/"+a.project),
	)
	req, err := http.NewRequest("POST", requestURI, body)
	if err != nil {
		return errors.Wrap(err, "prepare request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Private-Token", a.token)

	resp, err := a.HTTPClient().Do(req)
	if err != nil {
		return errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		log.Printf("created merge request of %s into %s", sourceBranchName, targetBranchName)

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "failed to read response body")
		}

		var response mergeRequestResponse
		err = json.Unmarshal(bodyBytes, &response)
		if err != nil {
			return errors.Wrap(err, "failed to decode response body")
		}

		log.Printf("visit the merge request at %s", response.URL)
	} else {
		return errors.Errorf("failed to create a merge_request: %s", resp.Status)
	}

	return nil
}

func (a *gitlabAdapter) GetCIBranch() (glow.Branch, error) {
	branchName := os.Getenv("CI_COMMIT_REF_NAME")
	branch, err := glow.BranchFromBranchName(branchName)
	return branch, err
}

func (a *gitlabAdapter) DetectCICDOrigin() (string, error) {
	gitProviderURL := a.endpoint
	gitUser := os.Getenv("CI_GIT_USER")
	gitToken := os.Getenv("CI_GIT_TOKEN")

	if gitUser != "" && gitToken != "" {
		endpointURL, err := url.Parse(a.endpoint)
		if err != nil {
			return "", errors.Wrap(err, "couldn't parse gitLab endpoint")
		}

		gitProviderURL = fmt.Sprintf(
			"%s://%s:%s@%s",
			endpointURL.Scheme, gitUser, gitToken, endpointURL.Host,
		)
	}

	return fmt.Sprintf(
		"%s/%s/%s.git",
		gitProviderURL, a.namespace, a.project,
	), nil
}
