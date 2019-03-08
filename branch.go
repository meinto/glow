package glow

import (
	"errors"
	"fmt"
	"strings"
)

// IBranch interface
type IBranch interface {
	CreationIsAllowedFrom(sourceBranch string) bool
	CanBeClosed() bool
	CanBePublished() bool
	CloseBranches(availableBranches []string) []string
	PublishBranch() string
	BranchName() string
	ShortBranchName() string
}

// GitService describes all actions which can performed with a branch
type GitService interface {
	CurrentBranch() (Branch, error)
	Create(b IBranch) error
	Checkout(b IBranch) error
}

// GitHostingService describes all actions which can performed
// with the git hosting service (gitlab etc.)
type GitHostingService interface {
	Close(b IBranch) error
	Publish(b IBranch) error
}

// Branch definition
type Branch struct {
	name string
}

// NewBranch creates a new branch definition
func NewBranch(name string) (Branch, error) {
	return Branch{name}, nil
}

// BranchName is a getter for the branch name
func (b Branch) BranchName() string {
	return b.name
}

// ShortBranchName is a getter for the short version of the branch name
func (b Branch) ShortBranchName() string {
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
		Branch{branchName},
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
		Branch{branchName},
	}, nil
}
