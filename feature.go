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
	if err != nil {
		return Feature{}, errors.Wrap(err, "error while creating feature definition")
	}
	return Feature{ab}, nil
}

// FeatureFromBranch extracts a feature definition from branch name
func FeatureFromBranch(branchName string) (Feature, error) {
	ab, err := AuthoredBranchFromBranchName(branchName)
	if err != nil {
		return Feature{}, errors.Wrap(err, "error while creating feature definition from branch name")
	}
	return Feature{ab}, nil
}

// CreationIsAllowed returns wheter branch is allowed to be created
// from given this source branch
func (f Feature) CreationIsAllowed(sourceBranch string) bool {
	if strings.Contains(sourceBranch, "develop") {
		return true
	}
	return false
}

// IsValid checks if the branch name is a valid
func (f Feature) IsValid() bool {
	if strings.Contains(f.name, "/feature/") {
		return true
	}
	return false
}

// FeatureService describes all actions which can performed with a feature
type FeatureService interface {
	Create(f Feature) error
	Close(f Feature) error
}
