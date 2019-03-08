package glow

import "fmt"

type fix struct {
	author     string
	name       string
	branchName string
}

func NewFix(author, name string) fix {
	branchName := fmt.Sprintf("refs/heads/fix/%s/%s", author, name)
	return fix{
		author,
		name,
		branchName,
	}
}
