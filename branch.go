package glow

import (
	"errors"
	"fmt"
	"strings"
)

// Branch interface
type Branch interface {
	CreationIsAllowedFrom(sourceBranch Branch) bool
	CanBeClosed() bool
	CanBePublished() bool
	CloseBranches(availableBranches []Branch) []Branch
	PublishBranch() Branch
	BranchName() string
	ShortBranchName() string
}

// Branch definition
type branch struct {
	name string
}

// NewBranch creates a new branch definition
func NewBranch(name string) Branch {
	if !strings.HasPrefix(name, "refs/heads/") {
		name = "refs/heads/" + name
	}
	return branch{name}
}

// CreationIsAllowedFrom returns wheter branch is allowed to be created
// from given this source branch
func (b branch) CreationIsAllowedFrom(sourceBranch Branch) bool {
	return false
}

// CanBeClosed checks if the branch name is a valid
func (b branch) CanBeClosed() bool {
	return false
}

// CanBePublished checks if the branch can be published directly to production
func (b branch) CanBePublished() bool {
	return false
}

// CloseBranches returns all branches which this branch have to be merged with
func (b branch) CloseBranches(availableBranches []Branch) []Branch {
	return []Branch{}
}

// PublishBranch returns the publish branch if available
func (b branch) PublishBranch() Branch {
	return nil
}

// BranchName is a getter for the branch name
func (b branch) BranchName() string {
	return b.name
}

// ShortBranchName is a getter for the short version of the branch name
func (b branch) ShortBranchName() string {
	return strings.TrimPrefix(b.name, "refs/heads/")
}

// AuthoredBranch definition
type AuthoredBranch interface {
	Branch
}

type authoredBranch struct {
	author string
	name   string
	Branch
}

// NewAuthoredBranch creates a new branch definition
func NewAuthoredBranch(branchTemplate, author, name string) (AuthoredBranch, error) {
	branchName := fmt.Sprintf(branchTemplate, author, name)
	branch := NewBranch(branchName)
	return authoredBranch{
		author,
		name,
		branch,
	}, nil
}

// AuthoredBranchFromBranchName extracts a feature definition from branch name
func AuthoredBranchFromBranchName(branchName string) (AuthoredBranch, error) {
	parts := strings.Split(branchName, "/")
	if len(parts) < 3 {
		return authoredBranch{}, errors.New("invalid branch name " + branchName)
	}
	branchType := parts[len(parts)-3]
	author := parts[len(parts)-2]
	name := parts[len(parts)-1]

	return NewAuthoredBranch(branchType+"/%s/%s", author, name)
}
