package gitprovider

import (
	"net/http"

	"github.com/imdario/mergo"
	"github.com/meinto/glow"
	"github.com/meinto/glow/cmd"
	"github.com/meinto/glow/git"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Service describes all actions which can performed
// with the git hosting git service (gitlab etc.)
type Service interface {
	GitService() git.Service
	HTTPClient() HttpClient
	SetHTTPClient(HttpClient)
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
	mainBranch string
	token      string
	gitService git.Service
	httpClient HttpClient
}

type Options struct {
	Endpoint   string
	Namespace  string
	Project    string
	MainBranch string
	Token      string
	ShouldLog  bool
	HttpClient HttpClient
}

var defaultOptions = Options{
	ShouldLog:  true,
	HttpClient: http.DefaultClient,
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
			options.MainBranch,
			options.Token,
			g,
			options.HttpClient,
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
			options.MainBranch,
			options.Token,
			g,
			options.HttpClient,
		},
	}

	if options.ShouldLog {
		NewLoggingService(s)
	}

	return s
}
