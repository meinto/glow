package glow

import (
	"strings"

	"github.com/pkg/errors"
)

// Hotfix definition
type Hotfix struct {
	AuthoredBranch
}

// NewHotfix creates a new hotfix definition
func NewHotfix(author, name string) (Hotfix, error) {
	ab, err := NewAuthoredBranch("refs/heads/hotfix/%s/%s", author, name)
	return Hotfix{ab}, errors.Wrap(err, "error while creating hotfix definition")
}

// HotfixFromBranch extracts a fix definition from branch name
func HotfixFromBranch(branchName string) (Hotfix, error) {
	if !strings.Contains(branchName, "/hotfix/") {
		return Hotfix{}, errors.New("no valid hotfix branch")
	}
	ab, err := AuthoredBranchFromBranchName(branchName)
	return Hotfix{ab}, errors.Wrap(err, "error while creating hotfix definition from branch name")
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
	develop, _ := NewBranch("develop")
	branches = append(branches, develop)
	return branches
}

// PublishBranch returns the publish branch if available
func (f Hotfix) PublishBranch() string {
	return "master"
}
