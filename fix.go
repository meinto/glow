package glow

import (
	"strings"

	"github.com/pkg/errors"
)

type Fix struct {
	AuthoredBranch
}

// NewFix creates a new fix definition
func NewFix(author, name string) (Fix, error) {
	ab, err := NewAuthoredBranch("refs/heads/fix/%s/%s", author, name)
	return Fix{ab}, errors.Wrap(err, "error while creating fix definition")
}

// FixFromBranch extracts a fix definition from branch name
func FixFromBranch(branchName string) (Fix, error) {
	if !strings.Contains(branchName, "/fix/") {
		return Fix{}, errors.New("no valid fix branch")
	}
	ab, err := AuthoredBranchFromBranchName(branchName)
	return Fix{ab}, errors.Wrap(err, "error while creating fix definition from branch name")
}

// CreationIsAllowedFrom returns wheter branch is allowed to be created
// from given this source branch
func (f Fix) CreationIsAllowedFrom(sourceBranch string) bool {
	return strings.Contains(sourceBranch, "release/v")
}

// CanBeClosed checks if the branch name is a valid
func (f Fix) CanBeClosed() bool {
	return true
}

// CloseBranches returns all branches which this branch have to be merged with
func (f Fix) CloseBranches(availableBranches []Branch) []Branch {
	branches := make([]Branch, 0)
	for _, b := range availableBranches {
		if strings.Contains(b.BranchName(), "/release/v") {
			branches = append(branches, b)
		}
	}
	return branches
}
