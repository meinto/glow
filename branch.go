package glow

import (
	"errors"
	"fmt"
	"strings"
)

// IBranch interface
type Branch interface {
	CreationIsAllowedFrom(sourceBranch string) bool
	CanBeClosed() bool
	CanBePublished() bool
	CloseBranches(availableBranches []Branch) []Branch
	PublishBranch() string
	BranchName() string
	ShortBranchName() string
}

// Branch definition
type branch struct {
	name string
}

// NewBranch creates a new branch definition
func NewBranch(name string) (Branch, error) {
	if strings.Contains(name, "/feature/") {
		return FeatureFromBranch(name)
	}
	if strings.Contains(name, "/fix/") {
		return FixFromBranch(name)
	}
	if strings.Contains(name, "/hotfix/") {
		return HotfixFromBranch(name)
	}
	if strings.Contains(name, "/release/v") {
		return ReleaseFromBranch(name)
	}
	return branch{name}, nil
}

// CreationIsAllowedFrom returns wheter branch is allowed to be created
// from given this source branch
func (b branch) CreationIsAllowedFrom(sourceBranch string) bool {
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
func (b branch) PublishBranch() string {
	return ""
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
type AuthoredBranch struct {
	author string
	name   string
	Branch
}

// NewAuthoredBranch creates a new branch definition
func NewAuthoredBranch(branchTemplate, author, name string) (AuthoredBranch, error) {
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
