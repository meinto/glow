package gitprovider

import (
	"os"

	"github.com/go-kit/kit/log"

	"github.com/meinto/glow"
	"github.com/meinto/glow/git"
)

// Service describes all actions which can performed
// with the git hosting git service (gitlab etc.)
type Service interface {
	Close(b glow.Branch) error
	Publish(b glow.Branch) error
}

type service struct {
	endpoint   string
	namespace  string
	project    string
	gitService git.Service
}

func NewGitlabService(endpoint, namespace, project, token string) Service {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	g := git.NewGoGitService()
	g = git.NewLoggingService(logger, g)

	return &gitlabAdapter{
		token,
		service{
			endpoint,
			namespace,
			project,
			g,
		},
	}
}

func NewGithubService(endpoint, namespace, project, clientID, clientSecret string) Service {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	g := git.NewGoGitService()
	g = git.NewLoggingService(logger, g)

	return &githubAdapter{
		clientID,
		clientSecret,
		service{
			endpoint,
			namespace,
			project,
			g,
		},
	}
}
