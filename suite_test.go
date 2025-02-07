package glow_test

import (
	"testing"

	. "github.com/meinto/glow"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

func TestGlow(t *testing.T) {
	viper.SetDefault("devBranch", "develop")
	viper.SetDefault("mainBranch", "master")
	RegisterFailHandler(Fail)
	RunSpecs(t, "Glow Suite")
}

const (
	MAIN_BRANCH     = "master"
	DEV_BRANCH      = "develop"
	RELEASE_BRANCH  = "release/v1.2.3"
	FEAUTURE_BRANCH = "feature/luke/falcon-shuttle"
	HOTFIX_BRANCH   = "hotfix/v0.0.1"
	FIX_BRANCH      = "fix/luke/falcon"
	ANOTHER_BRANCH  = "another-branch"
)

func MockBranchCollection() []Branch {
	b1 := NewBranch(MAIN_BRANCH)
	b2 := NewBranch(DEV_BRANCH)
	b3 := NewBranch(RELEASE_BRANCH)
	b4 := NewBranch(FEAUTURE_BRANCH)
	b5 := NewBranch(HOTFIX_BRANCH)
	b6 := NewBranch(FIX_BRANCH)
	b7 := NewBranch(ANOTHER_BRANCH)
	return []Branch{b1, b2, b3, b4, b5, b6, b7}
}
