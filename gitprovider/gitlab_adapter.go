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
	"github.com/pkg/errors"
)

type gitlabAdapter struct {
	service
}

func (a *gitlabAdapter) Close(b glow.Branch) error {
	if b.CanBeClosed() {
		branchList, err := a.gitService.BranchList()
		if err != nil {
			return errors.Wrap(err, "error getting branch list")
		}
		targets := b.CloseBranches(branchList)

		for _, t := range targets {
			err := a.createMergeRequest(b, t)
			if err != nil {
				return errors.Wrap(err, "error creating merge request")
			}
		}
		return nil
	}
	return errors.New("cannot be closed")
}

func (a *gitlabAdapter) Publish(b glow.Branch) error {
	if b.CanBePublished() {
		t := b.PublishBranch()
		return a.createMergeRequest(b, t)
	}
	return errors.New("cannot be published")
}

func (a *gitlabAdapter) createMergeRequest(source glow.Branch, target glow.Branch) error {
	type Payload struct {
		SourceBranch       string `json:"source_branch"`
		TargetBranch       string `json:"target_branch"`
		Title              string `json:"title"`
		RemoveSourceBranch bool   `json:"remove_source_branch"`
	}

	data := Payload{
		SourceBranch:       source.ShortBranchName(),
		TargetBranch:       target.ShortBranchName(),
		Title:              fmt.Sprintf("Merge %s in %s", source.ShortBranchName(), target.ShortBranchName()),
		RemoveSourceBranch: false,
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

	log.Printf("created merge request of %s into %s", source.ShortBranchName(), target.ShortBranchName())
	return nil
}

func (a *gitlabAdapter) GetCIBranch() (glow.Branch, error) {
	branchName := os.Getenv("CI_COMMIT_REF_NAME")
	branch, err := glow.NewBranch(branchName)
	return branch, errors.Wrap(err, "error get ci branch")
}
