package glow_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/meinto/glow"
	. "github.com/meinto/glow/testutil"
)

var _ = Describe("Release", func() {
	var branches []Branch

	BeforeEach(func() {
		f1, _ := NewRelease("1.2.3")
		f2, _ := ReleaseFromBranch("refs/heads/release/v1.2.3")
		branches = []Branch{f1, f2}
	})

	It("can be closed", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			Expect(branch.(Branch).CanBeClosed()).To(Equal(true))
		})
	})

	It("can be published", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			Expect(branch.(Branch).CanBePublished()).To(Equal(true))
		})
	})

	It("only closes on release branch", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			closeBanches := branch.(Branch).CloseBranches(MockBranchCollection())
			Expect(len(closeBanches)).To(Equal(1))
			Expect(closeBanches[0].ShortBranchName()).To(Equal(DEVELOP_BRANCH))
		})
	})

	It("is only allowed to create from develop branch", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			f := branch.(Branch)
			for _, testBranch := range MockBranchCollection() {
				testBranchName := testBranch.ShortBranchName()
				if testBranchName == DEVELOP_BRANCH {
					Expect(f.CreationIsAllowedFrom(testBranchName)).To(BeTrue())
				} else {
					Expect(f.CreationIsAllowedFrom(testBranchName)).To(BeFalse())
				}
			}
		})
	})

	It("publishes on the master branch", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			f := branch.(Branch)
			Expect(f.PublishBranch().ShortBranchName()).To(Equal(MASTER_BRANCH))
		})
	})

	// settings like default branch
	// ----------------------------
	It("has a branch name", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			branchName := branch.(Branch).BranchName()
			Expect(branchName).To(Equal("refs/heads/" + RELEASE_BRANCH))
		})
	})

	It("has a short branch name", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			branchName := branch.(Branch).ShortBranchName()
			Expect(branchName).To(Equal(RELEASE_BRANCH))
		})
	})
})
