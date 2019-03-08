package git

import "github.com/meinto/glow"

type Git struct {
	glow.GitService
}

func NewGit() Git {
	return Git{GoGitAdapter{}}
}

func NewNativeGit(gitPath string) Git {
	return Git{NativeGitAdapter{gitPath}}
}
