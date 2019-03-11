package gitprovider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/meinto/glow"
	"github.com/pkg/errors"
)

type githubAdapter struct {
	service
}

func (a *githubAdapter) Close(b glow.Branch) error {
	if b.CanBeClosed() {
		branchList, err := a.gitService.BranchList()
		if err != nil {
			return errors.Wrap(err, "error getting branch list")
		}
		targets := b.CloseBranches(branchList)

		for _, t := range targets {
			err := a.createPullRequest(b, t)
			if err != nil {
				return errors.Wrap(err, "error creating pull request")
			}
		}
		return nil
	}
	return errors.New("cannot be closed")
}

func (a *githubAdapter) Publish(b glow.Branch) error {
	if b.CanBePublished() {
		t := b.PublishBranch()
		return a.createPullRequest(b, t)
	}
	return errors.New("cannot be published")
}

func (a *githubAdapter) createPullRequest(source glow.Branch, target glow.Branch) error {
	type Payload struct {
		Head                string `json:"head"`
		Base                string `json:"base"`
		Title               string `json:"title"`
		MaintainerCanModify bool   `json:"maintainer_can_modify"`
	}

	data := Payload{
		Head:                source.ShortBranchName(),
		Base:                target.ShortBranchName(),
		Title:               fmt.Sprintf("Pull %s in %s", source.ShortBranchName(), target.ShortBranchName()),
		MaintainerCanModify: true,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "marshal json")
	}

	body := bytes.NewReader(payloadBytes)

	requestURI := fmt.Sprintf(
		"%s/repos/%s/%s/pulls",
		a.endpoint,
		a.namespace,
		a.project,
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

	log.Printf("created pull request of %s into %s", source.ShortBranchName(), target.ShortBranchName())
	return nil
}
