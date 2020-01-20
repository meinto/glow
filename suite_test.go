package glow_test

import (
	"testing"

	. "github.com/meinto/glow"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGlow(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Glow Suite")
}

const (
	MASTER_BRANCH   = "master"
	DEVELOP_BRANCH  = "develop"
	RELEASE_BRANCH  = "release/v1.2.3"
	FEAUTURE_BRANCH = "feature/luke/falcon"
	HOTFIX_BRANCH   = "hotfix/v0.0.1"
	FIX_BRANCH      = "fix/luke/fix-falcon"
	ANOTHER_BRANCH  = "another-branch"
)

func MockBranchCollection() []Branch {
	b1, _ := NewBranch(MASTER_BRANCH)
	b2, _ := NewBranch(DEVELOP_BRANCH)
	b3, _ := NewBranch(RELEASE_BRANCH)
	b4, _ := NewBranch(FEAUTURE_BRANCH)
	b5, _ := NewBranch(HOTFIX_BRANCH)
	b6, _ := NewBranch(FIX_BRANCH)
	b7, _ := NewBranch(ANOTHER_BRANCH)
	return []Branch{b1, b2, b3, b4, b5, b6, b7}
}
