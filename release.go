package glow

import "fmt"

type release struct {
	version    string
	branchName string
}

func NewRelease(version, name string) release {
	branchName := fmt.Sprintf("refs/heads/release/v%s", version)
	return release{
		version,
		branchName,
	}
}
