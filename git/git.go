package git

import "github.com/meinto/glow"

type Git struct {
	adapter glow.GitService
}

func NewGit(a glow.GitService) Git {
	return Git{a}
}
