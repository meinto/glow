package glow

import (
	"fmt"
	"strings"

	l "github.com/meinto/glow/logging"
	"github.com/pkg/errors"
)

// Release definition
type release struct {
	version string
	Branch
}

// NewRelease creates a new release definition
func NewRelease(version string) (b Branch, err error) {
	l.Log().Info(l.Fields{"version": version})
	defer func() {
		l.Log().
			Info(l.Fields{"branch": b}).
			Error(err)
	}()
	branchName := fmt.Sprintf(BRANCH_NAME_PREFIX+"release/v%s", version)
	b = NewBranch(branchName)
	return release{version, b}, nil
}

// ReleaseFromBranch extracts a release definition from branch name
func ReleaseFromBranch(branchName string) (b Branch, err error) {
	l.Log().Info(l.Fields{"branchName": branchName})
	defer func() {
		l.Log().
			Info(l.Fields{"branch": b}).
			Error(err)
	}()
	if !strings.Contains(branchName, "/release/v") {
		return release{}, errors.New("no valid release branch")
	}
	b = NewBranch(branchName)
	parts := strings.Split(branchName, "/")
	if len(parts) < 1 {
		return release{}, errors.New("invalid branch name " + branchName)
	}
	version := parts[len(parts)-1]
	version = strings.TrimPrefix(version, "v")

	return release{version, b}, nil
}

// CreationIsAllowedFrom returns wheter branch is allowed to be created
// from given this source branch
func (f release) CreationIsAllowedFrom(sourceBranch Branch) bool {
	if strings.Contains(sourceBranch.ShortBranchName(), "develop") {
		return true
	}
	return false
}

// CanBeClosed checks if the branch name is a valid
func (f release) CanBeClosed() bool {
	return true
}

// CanBePublished checks if the branch can be published directly to production
func (f release) CanBePublished() bool {
	return true
}

// CloseBranches returns all branches which this branch have to be merged with
func (f release) CloseBranches(availableBranches []Branch) []Branch {
	return []Branch{
		NewBranch("develop"),
	}
}

// PublishBranch returns the publish branch if available
func (f release) PublishBranch() Branch {
	return NewBranch("master")
}
