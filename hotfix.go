package glow

import (
	"fmt"
	"regexp"
	"strings"

	l "github.com/meinto/glow/logging"
	"github.com/pkg/errors"
)

// Hotfix definition
type hotfix struct {
	version string
	Branch
}

// NewHotfix creates a new hotfix definition
func NewHotfix(version string) (b Branch, err error) {
	l.Log().Info(l.Fields{"version": version})
	defer func() {
		l.Log().
			Info(l.Fields{"branch": b}).
			Error(err)
	}()
	branchName := fmt.Sprintf(BRANCH_NAME_PREFIX+"hotfix/v%s", version)
	b = NewBranch(branchName)
	return hotfix{version, b}, nil
}

// HotfixFromBranch extracts a fix definition from branch name
func HotfixFromBranch(branchName string) (b Branch, err error) {
	l.Log().Info(l.Fields{"branchName": branchName})
	defer func() {
		l.Log().
			Info(l.Fields{"branch": b}).
			Error(err)
	}()
	matched, err := regexp.Match(HOTFIX_BRANCH_PATTERN, []byte(branchName))
	if !matched || err != nil {
		return hotfix{}, errors.New("no valid hotfix branch")
	}
	b = NewBranch(branchName)
	parts := strings.Split(branchName, "/")
	if len(parts) < 1 {
		return hotfix{}, errors.New("invalid branch name " + branchName)
	}
	version := parts[len(parts)-1]
	version = strings.TrimPrefix(version, "v")

	return hotfix{version, b}, nil
}

// CreationIsAllowedFrom returns wheter branch is allowed to be created
// from given this source branch
func (f hotfix) CreationIsAllowedFrom(sourceBranch Branch) bool {
	if strings.Contains(sourceBranch.ShortBranchName(), "master") {
		return true
	}
	return false
}

// CanBeClosed checks if the branch name is a valid
func (f hotfix) CanBeClosed() bool {
	return true
}

// CanBePublished checks if the branch can be published directly to production
func (f hotfix) CanBePublished() bool {
	return true
}

// CloseBranches returns all branches which this branch have to be merged with
func (f hotfix) CloseBranches(availableBranches []Branch) []Branch {
	branches := make([]Branch, 0)
	for _, b := range availableBranches {
		if strings.Contains(b.BranchName(), "/release/v") {
			branches = append(branches, b)
		}
	}
	branches = append(branches, NewBranch("develop"))
	return branches
}

// PublishBranch returns the publish branch if available
func (f hotfix) PublishBranch() Branch {
	return NewBranch("master")
}
