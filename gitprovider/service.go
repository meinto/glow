package gitprovider

import (
	"io"
	"net/http"

	"github.com/imdario/mergo"
	"github.com/meinto/glow"
	"github.com/meinto/glow/cmd"
	"github.com/meinto/glow/git"
)

// Service describes all actions which can performed
// with the git hosting git service (gitlab etc.)
type Service interface {
	GitService() git.Service
	NewRequest(method, url string, body io.Reader) (*http.Request, error)
	SetGitService(git.Service)
	Close(b glow.Branch) error
	Publish(b glow.Branch) error
	DetectCICDOrigin() (string, error)
	GetCIBranch() (glow.Branch, error)
}

type service struct {
	endpoint   string
	namespace  string
	project    string
	token      string
	gitService git.Service
}

type Options struct {
	Endpoint  string
	Namespace string
	Project   string
	Token     string
	ShouldLog bool
}

var defaultOptions = Options{
	ShouldLog: true,
}

func NewGitlabService(options Options) Service {
	mergo.Merge(&options, defaultOptions)
	exec := cmd.NewCmdExecutor("/bin/bash")
	g := git.NewNativeService(git.Options{
		CmdExecutor: exec,
	})

	s := &gitlabAdapter{
		service{
			options.Endpoint,
			options.Namespace,
			options.Project,
			options.Token,
			g,
		},
	}

	if options.ShouldLog {
		NewLoggingService(s)
	}

	return s
}

func NewGithubService(options Options) Service {
	mergo.Merge(&options, defaultOptions)
	exec := cmd.NewCmdExecutor("/bin/bash")
	g := git.NewNativeService(git.Options{
		CmdExecutor: exec,
	})

	s := &githubAdapter{
		service{
			options.Endpoint,
			options.Namespace,
			options.Project,
			options.Token,
			g,
		},
	}

	if options.ShouldLog {
		NewLoggingService(s)
	}

	return s
}
