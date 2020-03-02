package glow_test

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/meinto/glow"
	. "github.com/meinto/glow/testutil"
)

var _ = Describe("Hotfix", func() {
	var branches []Branch

	BeforeEach(func() {
		f1, _ := NewHotfix("0.0.1")
		f2, _ := HotfixFromBranch(BRANCH_NAME_PREFIX + "hotfix/v0.0.1")
		f3, _ := BranchFromBranchName(BRANCH_NAME_PREFIX + "hotfix/v0.0.1")
		f4, _ := BranchFromBranchName("hotfix/v0.0.1")
		branches = []Branch{f1, f2, f3, f4}
	})

	It("is of type hotfix branch", func() {
		f, _ := NewHotfix("a")
		rf := reflect.ValueOf(f)
		ForEachTestSet(branches, func(branch interface{}) {
			r := reflect.ValueOf(branch)
			Expect(r.Type().AssignableTo(rf.Type())).To(BeTrue())
		})
	})

	It("can be closed", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			Expect(branch.(Branch).CanBeClosed()).To(Equal(true))
		})
	})

	It("closes on release branches & develop branch", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			closeBanches := branch.(Branch).CloseBranches(MockBranchCollection())
			Expect(len(closeBanches)).To(Equal(2))
			Expect(closeBanches[0].ShortBranchName()).To(Equal(RELEASE_BRANCH))
			Expect(closeBanches[1].ShortBranchName()).To(Equal(DEVELOP_BRANCH))
		})
	})

	It("is only allowed to create from master branch", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			f := branch.(Branch)
			for _, testBranch := range MockBranchCollection() {
				testBranchName := testBranch.ShortBranchName()
				if testBranchName == MASTER_BRANCH {
					Expect(f.CreationIsAllowedFrom(testBranch)).To(BeTrue())
				} else {
					Expect(f.CreationIsAllowedFrom(testBranch)).To(BeFalse())
				}
			}
		})
	})

	It("can be published", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			Expect(branch.(Branch).CanBePublished()).To(BeTrue())
		})
	})

	It("can be published on master", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			publishBranch := branch.(Branch).PublishBranch()
			Expect(publishBranch.ShortBranchName()).To(Equal(MASTER_BRANCH))
		})
	})

	// settings like default branch
	// ----------------------------
	It("has a branch name", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			branchName := branch.(Branch).BranchName()
			Expect(branchName).To(Equal(BRANCH_NAME_PREFIX + HOTFIX_BRANCH))
		})
	})

	It("has a short branch name", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			branchName := branch.(Branch).ShortBranchName()
			Expect(branchName).To(Equal(HOTFIX_BRANCH))
		})
	})
})
