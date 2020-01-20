package glow

import (
	"strings"

	"github.com/pkg/errors"
)

type fix struct {
	AuthoredBranch
}

// NewFix creates a new fix definition
func NewFix(author, name string) (AuthoredBranch, error) {
	ab, err := NewAuthoredBranch(AUTHORED_BRANCH_TYPE_FIX, author, name)
	return fix{ab}, errors.Wrap(err, "error while creating fix definition")
}

// FixFromBranch extracts a fix definition from branch name
func FixFromBranch(branchName string) (AuthoredBranch, error) {
	if !strings.Contains(branchName, "/fix/") {
		return fix{}, errors.New("no valid fix branch")
	}
	ab, err := AuthoredBranchFromBranchName(branchName)
	return fix{ab}, errors.Wrap(err, "error while creating fix definition from branch name")
}

// CreationIsAllowedFrom returns wheter branch is allowed to be created
// from given this source branch
func (f fix) CreationIsAllowedFrom(sourceBranch Branch) bool {
	return strings.Contains(sourceBranch.ShortBranchName(), "release/v")
}

// CanBeClosed checks if the branch name is a valid
func (f fix) CanBeClosed() bool {
	return true
}

// CloseBranches returns all branches which this branch have to be merged with
func (f fix) CloseBranches(availableBranches []Branch) []Branch {
	branches := make([]Branch, 0)
	for _, b := range availableBranches {
		if strings.Contains(b.BranchName(), "/release/v") {
			branches = append(branches, b)
		}
	}
	return branches
}
