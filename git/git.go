package git

import "github.com/meinto/glow"

type Git struct {
	gogit  glow.GitService
	native glow.GitService
}

func NewGit() Git {
	return Git{
		GoGitAdapter{},
		NativeGitAdapter{},
	}
}
