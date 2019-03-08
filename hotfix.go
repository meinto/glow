package glow

import "fmt"

type hotfix struct {
	author     string
	name       string
	branchName string
}

func NewHotfix(author, name string) hotfix {
	branchName := fmt.Sprintf("refs/heads/hotfix/%s/%s", author, name)
	return hotfix{
		author,
		name,
		branchName,
	}
}
