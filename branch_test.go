package glow_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/meinto/glow"
	. "github.com/meinto/glow/testutil"
)

var _ = Describe("Branch", func() {
	var branches []Branch

	BeforeEach(func() {
		b1 := NewBranch("fix/author/feature")
		b3, _ := NewAuthoredBranch("fix", "author", "feature")
		b4, _ := AuthoredBranchFromBranchName("fix/author/feature")
		branches = []Branch{b1, b3, b4}
	})

	It("is not allowed to be created from another branch", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			Expect(branch.(Branch).CreationIsAllowedFrom(NewBranch("another-branch"))).To(BeFalse())
		})
	})

	It("cannot be closed", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			Expect(branch.(Branch).CanBeClosed()).To(BeFalse())
		})
	})

	It("cannot be published", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			Expect(branch.(Branch).CanBePublished()).To(BeFalse())
		})
	})

	It("has no close branches", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			closeBranches := branch.(Branch).CloseBranches(MockBranchCollection())
			Expect(len(closeBranches)).To(Equal(0))
		})
	})

	It("has no publish branch", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			publishBranch := branch.(Branch).PublishBranch()
			Expect(publishBranch).To(BeNil())
		})
	})

	It("has a branch name", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			branchName := branch.(Branch).BranchName()
			Expect(branchName).To(Equal(BRANCH_NAME_PREFIX + "fix/author/feature"))
		})
	})

	It("has a short branch name", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			branchName := branch.(Branch).ShortBranchName()
			Expect(branchName).To(Equal("fix/author/feature"))
		})
	})
})
