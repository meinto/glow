package gitprovider

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func (s *gitlabAdapter) SetGitService(gs git.Service) {
	s.gitService = gs
}

func (a *gitlabAdapter) Close(b glow.Branch) error {
	_, _, remoteBranchExists := a.gitService.RemoteBranchExists(b.ShortBranchName())
	if b.CanBeClosed() && remoteBranchExists == nil {
		branchList, _, _, err := a.gitService.BranchList()
		if err != nil {
			return errors.Wrap(err, "error getting branch list")
		}
		targets := b.CloseBranches(branchList)

		for _, t := range targets {
			err := a.createMergeRequest(b, t, true)
			if err != nil {
				return errors.Wrap(err, "error creating merge request")
			}
		}
		return nil
	}
	return errors.Wrap(remoteBranchExists, "cannot be closed")
}

func (a *gitlabAdapter) Publish(b glow.Branch) error {
	_, _, remoteBranchExists := a.gitService.RemoteBranchExists(b.ShortBranchName())
	if b.CanBePublished() && remoteBranchExists == nil {
		t := b.PublishBranch()
		return a.createMergeRequest(b, t, false)
	}
	return errors.Wrap(remoteBranchExists, "cannot be published")
}

func (a *gitlabAdapter) createMergeRequest(source glow.Branch, target glow.Branch, removeSourceBranch bool) error {
	type Payload struct {
		SourceBranch       string `json:"source_branch"`
		TargetBranch       string `json:"target_branch"`
		Title              string `json:"title"`
		RemoveSourceBranch bool   `json:"remove_source_branch"`
		Squash             bool   `json:"squash"`
	}

	sourceBranchName := source.ShortBranchName()
	targetBranchName := target.ShortBranchName()

	data := Payload{
		SourceBranch:       sourceBranchName,
		TargetBranch:       targetBranchName,
		Title:              fmt.Sprintf("Merge %s in %s", sourceBranchName, targetBranchName),
		RemoveSourceBranch: removeSourceBranch,
		Squash:             viper.GetBool("mergeRequest.squashCommits"),
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

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	log.Printf("created merge request of %s into %s", sourceBranchName, targetBranchName)
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
