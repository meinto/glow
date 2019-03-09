package githost

import (
	"github.com/meinto/glow"
	"github.com/pkg/errors"
)

type gitlabAdapter struct {
	service
}

func (a *gitlabAdapter) Close(b glow.Branch) error {
	if b.CanBeClosed() {
		// a.createMergeRequest(b)
	}
	return errors.New("not implemented yet")
}

func (a *gitlabAdapter) Publish(b glow.Branch) error {
	return errors.New("not implemented yet")
}

// func (a *gitlabAdapter) createMergeRequest(b glow.Branch) error {
// 	type Payload struct {
// 		SourceBranch       string `json:"source_branch"`
// 		TargetBranch       string `json:"target_branch"`
// 		Title              string `json:"title"`
// 		RemoveSourceBranch bool   `json:"remove_source_branch"`
// 	}

// 	currentBranch := a.gitService.CurrentBranch()
// 	targets := b.CloseBranches()

// 	data := Payload{
// 		SourceBranch:       b.ShortBranchName(),
// 		TargetBranch:       target.ShortBranchName(),
// 		Title:              fmt.Sprintf("Merge %s in %s", source.ShortBranchName(), target.ShortBranchName()),
// 		RemoveSourceBranch: false,
// 	}
// 	payloadBytes, err := json.Marshal(data)
// 	CheckForError(err, "mergeRequestCmd marshal json")

// 	body := bytes.NewReader(payloadBytes)

// 	requestURI := fmt.Sprintf(
// 		"%s/api/v4/projects/%s%s%s/merge_requests",
// 		a.endpoint,
// 		a.namespace,
// 		"%2F",
// 		a.project,
// 	)
// 	req, err := http.NewRequest("POST", requestURI, body)
// 	CheckForError(err, "mergeRequestCmd prepare request")

// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Private-Token", a.token)

// 	resp, err := http.DefaultClient.Do(req)
// 	CheckForError(err, "mergeRequestCmd do request")
// 	defer resp.Body.Close()

// 	log.Printf("created merge request of %s into %s", source.ShortBranchName(), target.ShortBranchName())
// }
