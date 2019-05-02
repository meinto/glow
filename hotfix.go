package glow

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Hotfix definition
type Hotfix struct {
	version string
	Branch
}

// NewHotfix creates a new hotfix definition
func NewHotfix(version string) (Hotfix, error) {
	branchName := fmt.Sprintf("refs/heads/hotfix/v%s", version)
	b := NewPlainBranch(branchName)
	return Hotfix{version, b}, nil
}

// HotfixFromBranch extracts a fix definition from branch name
func HotfixFromBranch(branchName string) (Hotfix, error) {
	if !strings.Contains(branchName, "/hotfix/v") {
		return Hotfix{}, errors.New("no valid hotfix branch")
	}
	b := NewPlainBranch(branchName)
	parts := strings.Split(branchName, "/")
	if len(parts) < 1 {
		return Hotfix{}, errors.New("invalid branch name " + branchName)
	}
	version := parts[len(parts)-1]
	version = strings.TrimPrefix(version, "v")

	return Hotfix{version, b}, nil
}

// CreationIsAllowedFrom returns wheter branch is allowed to be created
// from given this source branch
func (f Hotfix) CreationIsAllowedFrom(sourceBranch string) bool {
	if strings.Contains(sourceBranch, "master") {
		return true
	}
	return false
}

// CanBeClosed checks if the branch name is a valid
func (f Hotfix) CanBeClosed() bool {
	return true
}

// CanBePublished checks if the branch can be published directly to production
func (f Hotfix) CanBePublished() bool {
	return true
}

// CloseBranches returns all branches which this branch have to be merged with
func (f Hotfix) CloseBranches(availableBranches []Branch) []Branch {
	branches := make([]Branch, 0)
	for _, b := range availableBranches {
		if strings.Contains(b.BranchName(), "/release/v") {
			branches = append(branches, b)
		}
	}
	branches = append(branches, NewPlainBranch("develop"))
	return branches
}

// PublishBranch returns the publish branch if available
func (f Hotfix) PublishBranch() Branch {
	return NewPlainBranch("master")
}
