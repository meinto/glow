package glow

import (
	"errors"
	"fmt"
	"strings"
)

// Branch interface
type Branch interface {
	CreationIsAllowed(sourceBranch string) bool
	IsValid() bool
}

// Branch definition
type branch struct {
	name string
}

// AuthoredBranch definition
type AuthoredBranch struct {
	author string
	name   string
	branch
}

// NewAuthoredBranch creates a new branch definition
func NewAuthoredBranch(author, name, branchTemplate string) (AuthoredBranch, error) {
	branchName := fmt.Sprintf(branchTemplate, author, name)
	return AuthoredBranch{
		author,
		name,
		branch{branchName},
	}, nil
}

// AuthoredBranchFromBranchName extracts a feature definition from branch name
func AuthoredBranchFromBranchName(branchName string) (AuthoredBranch, error) {
	parts := strings.Split(branchName, "/")
	if len(parts) < 2 {
		return AuthoredBranch{}, errors.New("invalid branch name " + branchName)
	}
	author := parts[len(parts)-2]
	name := parts[len(parts)-1]

	return AuthoredBranch{
		author,
		name,
		branch{branchName},
	}, nil
}
