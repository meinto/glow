package glow

import (
	"errors"
	"fmt"
	"strings"

	l "github.com/meinto/glow/logging"
)

// Branch interface
type Branch interface {
	CreationIsAllowedFrom(sourceBranch Branch) bool
	CanBeClosed() bool
	CanBePublished() bool
	CloseBranches(availableBranches []Branch) []Branch
	PublishBranch() Branch
	BranchName() string
	ShortBranchName() string
}

// Branch definition
type branch struct {
	name string
}

const BRANCH_NAME_PREFIX = "refs/heads/"

// NewBranch creates a new branch definition
func NewBranch(name string) Branch {
	l.Log().Debug(l.Fields{"name": name})
	if !strings.HasPrefix(name, BRANCH_NAME_PREFIX) {
		name = BRANCH_NAME_PREFIX + name
	}
	return NewBranchLoggingService(branch{name})
}

func BranchFromBranchName(name string) (b Branch, err error) {
	l.Log().Debug(l.Fields{"name": name})
	defer func() {
		l.Log().
			Debug(l.Fields{"branch": b}).
			Error(err)
	}()
	if strings.Contains(name, "/feature/") {
		return FeatureFromBranch(name)
	}
	if strings.Contains(name, "/fix/") {
		return FixFromBranch(name)
	}
	if strings.Contains(name, "/hotfix/") {
		return HotfixFromBranch(name)
	}
	if strings.Contains(name, "/release/v") {
		return ReleaseFromBranch(name)
	}
	return NewBranch(name), nil
}

// CreationIsAllowedFrom returns wheter branch is allowed to be created
// from given this source branch
func (b branch) CreationIsAllowedFrom(sourceBranch Branch) bool {
	return false
}

// CanBeClosed checks if the branch name is a valid
func (b branch) CanBeClosed() bool {
	return false
}

// CanBePublished checks if the branch can be published directly to production
func (b branch) CanBePublished() bool {
	return false
}

// CloseBranches returns all branches which this branch have to be merged with
func (b branch) CloseBranches(availableBranches []Branch) []Branch {
	return []Branch{}
}

// PublishBranch returns the publish branch if available
func (b branch) PublishBranch() Branch {
	return nil
}

// BranchName is a getter for the branch name
func (b branch) BranchName() string {
	return b.name
}

// ShortBranchName is a getter for the short version of the branch name
func (b branch) ShortBranchName() string {
	return strings.TrimPrefix(b.name, BRANCH_NAME_PREFIX)
}

// AuthoredBranch definition
type AuthoredBranch interface {
	Branch
}

type authoredBranch struct {
	author      string
	featureName string
	Branch
}

const (
	AUTHORED_BRANCH_TYPE_FEATURE = "feature"
	AUTHORED_BRANCH_TYPE_FIX     = "fix"
)

// NewAuthoredBranch creates a new branch definition
func NewAuthoredBranch(branchType, author, featureName string) (AuthoredBranch, error) {
	if branchType != AUTHORED_BRANCH_TYPE_FEATURE && branchType != AUTHORED_BRANCH_TYPE_FIX {
		return nil, fmt.Errorf("branch type '%s' is not valid for an authored branch", branchType)
	}
	branchName := fmt.Sprintf("%s/%s/%s", branchType, author, featureName)
	branch := NewBranch(branchName)
	return authoredBranch{
		author,
		featureName,
		branch,
	}, nil
}

// AuthoredBranchFromBranchName extracts a feature definition from branch name
func AuthoredBranchFromBranchName(branchName string) (AuthoredBranch, error) {
	parts := strings.Split(branchName, "/")
	if len(parts) < 3 {
		return authoredBranch{}, errors.New("invalid branch name " + branchName)
	}
	branchType := parts[len(parts)-3]
	author := parts[len(parts)-2]
	featureName := parts[len(parts)-1]

	return NewAuthoredBranch(branchType, author, featureName)
}
