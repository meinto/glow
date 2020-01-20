package glow

import (
	"strings"

	"github.com/pkg/errors"
)

// Feature definition
type Feature struct {
	AuthoredBranch
}

// NewFeature creates a new feature definition
func NewFeature(author, name string) (Feature, error) {
	ab, err := NewAuthoredBranch("refs/heads/feature/%s/%s", author, name)
	return Feature{ab}, errors.Wrap(err, "error while creating feature definition")
}

// FeatureFromBranch extracts a feature definition from branch name
func FeatureFromBranch(branchName string) (Feature, error) {
	if !strings.Contains(branchName, "/feature/") {
		return Feature{}, errors.New("no valid feature branch")
	}
	ab, err := AuthoredBranchFromBranchName(branchName)
	return Feature{ab}, errors.Wrap(err, "error while creating feature definition from branch name")
}

// CreationIsAllowedFrom returns wheter branch is allowed to be created
// from given this source branch
func (f Feature) CreationIsAllowedFrom(sourceBranch string) bool {
	return strings.Contains(sourceBranch, "develop")
}

// CanBeClosed checks if the branch name is a valid
func (f Feature) CanBeClosed() bool {
	return true
}

// CloseBranches returns all branches which this branch have to be merged with
func (f Feature) CloseBranches(availableBranches []Branch) []Branch {
	return []Branch{
		NewPlainBranch("develop"),
	}
}
