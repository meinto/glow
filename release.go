package glow

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Release definition
type Release struct {
	version string
	Branch
}

// NewRelease creates a new release definition
func NewRelease(version string) (Release, error) {
	branchName := fmt.Sprintf("refs/heads/release/v%s", version)
	b, err := NewBranch(branchName)
	return Release{version, b}, errors.Wrap(err, "error while creating release definition")
}

// ReleaseFromBranch extracts a release definition from branch name
func ReleaseFromBranch(branchName string) (Release, error) {
	if !strings.Contains(branchName, "/release/v") {
		return Release{}, errors.New("no valid release branch")
	}
	b, err := NewBranch(branchName)
	if err != nil {
		return Release{}, errors.Wrap(err, "error while creating release definition from branch name")
	}
	parts := strings.Split(branchName, "/")
	if len(parts) < 1 {
		return Release{}, errors.New("invalid branch name " + branchName)
	}
	version := parts[len(parts)-1]
	version = strings.TrimPrefix(version, "v")

	return Release{version, b}, nil
}

// CreationIsAllowedFrom returns wheter branch is allowed to be created
// from given this source branch
func (f Release) CreationIsAllowedFrom(sourceBranch string) bool {
	if strings.Contains(sourceBranch, "develop") {
		return true
	}
	return false
}

// CanBeClosed checks if the branch name is a valid
func (f Release) CanBeClosed() bool {
	return true
}

// CanBePublished checks if the branch can be published directly to production
func (f Release) CanBePublished() bool {
	return true
}

// CloseBranches returns all branches which this branch have to be merged with
func (f Release) CloseBranches(availableBranches []Branch) []Branch {
	develop, _ := NewBranch("develop")
	return []Branch{
		develop,
	}
}

// PublishBranch returns the publish branch if available
func (f Release) PublishBranch() string {
	return "master"
}
