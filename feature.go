package glow

import (
	"regexp"
	"strings"

	l "github.com/meinto/glow/logging"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Feature definition
type feature struct {
	AuthoredBranch
}

// NewFeature creates a new feature definition
func NewFeature(author, name string) (b AuthoredBranch, err error) {
	l.Log().Info(l.Fields{"author": author, "name": name})
	defer func() {
		l.Log().
			Info(l.Fields{"branch": b}).
			Error(err)
	}()
	ab, err := NewAuthoredBranch(AUTHORED_BRANCH_TYPE_FEATURE, author, name)
	return feature{ab}, errors.Wrap(err, "error while creating feature definition")
}

// FeatureFromBranch extracts a feature definition from branch name
func FeatureFromBranch(branchName string) (b AuthoredBranch, err error) {
	l.Log().Info(l.Fields{"branchName": branchName})
	defer func() {
		l.Log().
			Info(l.Fields{"branch": b}).
			Error(err)
	}()
	matched, err := regexp.Match(FEATURE_BRANCH_PATTERN, []byte(branchName))
	if !matched || err != nil {
		return feature{}, errors.New("no valid feature branch")
	}
	ab, err := AuthoredBranchFromBranchName(branchName)
	return feature{ab}, errors.Wrap(err, "error while creating feature definition from branch name")
}

// CreationIsAllowedFrom returns wheter branch is allowed to be created
// from given this source branch
func (f feature) CreationIsAllowedFrom(sourceBranch Branch) bool {
	matched, err := regexp.Match(FEATURE_BRANCH_PATTERN, []byte(sourceBranch.BranchName()))
	devBranch := viper.GetString("devBranch")
	return strings.Contains(sourceBranch.ShortBranchName(), devBranch) || (matched && err == nil)
}

// CanBeClosed checks if the branch name is a valid
func (f feature) CanBeClosed() bool {
	return true
}

// CloseBranches returns all branches which this branch have to be merged with
func (f feature) CloseBranches(availableBranches []Branch) []Branch {
	devBranch := viper.GetString("devBranch")
	return []Branch{
		NewBranch(devBranch),
	}
}
