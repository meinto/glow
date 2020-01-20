package gitprovider

import (
	"github.com/meinto/glow"
	"github.com/meinto/glow/cmd"
	"github.com/meinto/glow/git"
)

// Service describes all actions which can performed
// with the git hosting git service (gitlab etc.)
type Service interface {
	Close(b glow.Branch) error
	Publish(b glow.Branch) error
	DetectCICDOrigin() (string, error)
	GetCIBranch() glow.Branch
}

type service struct {
	endpoint   string
	namespace  string
	project    string
	token      string
	gitService git.Service
}

func NewGitlabService(endpoint, namespace, project, token string) Service {
	exec := cmd.NewCmdExecutor("/bin/bash")
	g := git.NewNativeService(exec)
	g = git.NewLoggingService(g)

	return &gitlabAdapter{
		service{
			endpoint,
			namespace,
			project,
			token,
			g,
		},
	}
}

func NewGithubService(endpoint, namespace, project, token string) Service {
	exec := cmd.NewCmdExecutor("/bin/bash")
	g := git.NewNativeService(exec)
	g = git.NewLoggingService(g)

	return &githubAdapter{
		service{
			endpoint,
			namespace,
			project,
			token,
			g,
		},
	}
}
