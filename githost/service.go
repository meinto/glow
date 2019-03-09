package githost

import (
	"github.com/meinto/glow"
)

// Service describes all actions which can performed
// with the git hosting git service (gitlab etc.)
type Service interface {
	Close(b glow.IBranch) error
	Publish(b glow.IBranch) error
}

type service struct {
	endpoint  string
	namespace string
	project   string
	token     string
}

func NewGitlabService(endpoint, namespace, project, token string) Service {
	return &gitlabAdapter{
		service{
			endpoint,
			namespace,
			project,
			token,
		},
	}
}

// type Payload struct {
// 	SourceBranch       string `json:"source_branch"`
// 	TargetBranch       string `json:"target_branch"`
// 	Title              string `json:"title"`
// 	RemoveSourceBranch bool   `json:"remove_source_branch"`
// }

// data := Payload{
// 	SourceBranch:       source,
// 	TargetBranch:       target,
// 	Title:              fmt.Sprintf("Merge %s in %s", source, target),
// 	RemoveSourceBranch: false,
// }
// payloadBytes, err := json.Marshal(data)
// CheckForError(err, "mergeRequestCmd marshal json")

// body := bytes.NewReader(payloadBytes)

// requestURI := fmt.Sprintf(
// 	"%s/api/v4/projects/%s%s%s/merge_requests",
// 	viper.GetString("gitlabEndpoint"),
// 	viper.GetString("projectNamespace"),
// 	"%2F",
// 	viper.GetString("projectName"),
// )
// req, err := http.NewRequest("POST", requestURI, body)
// CheckForError(err, "mergeRequestCmd prepare request")

// req.Header.Set("Content-Type", "application/json")
// req.Header.Set("Private-Token", viper.GetString("gitlabCIToken"))

// resp, err := http.DefaultClient.Do(req)
// CheckForError(err, "mergeRequestCmd do request")
// defer resp.Body.Close()

// log.Printf("created merge request of %s into %s", source, target)
