package glow_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/meinto/glow"
	. "github.com/meinto/glow/testutil"
)

var _ = Describe("Feature", func() {
	var features []AuthoredBranch

	BeforeEach(func() {
		f1, _ := NewFeature("luke", "falcon")
		f2, _ := FeatureFromBranch("refs/heads/feature/luke/falcon")
		features = []AuthoredBranch{f1, f2}
	})

	It("can be closed", func() {
		ForEachTestSet(features, func(feature interface{}) {
			Expect(feature.(AuthoredBranch).CanBeClosed()).To(Equal(true))
		})
	})

	It("only closes on develop", func() {
		ForEachTestSet(features, func(feature interface{}) {
			closeBanches := feature.(AuthoredBranch).CloseBranches([]Branch{})
			Expect(len(closeBanches)).To(Equal(1))
			Expect(closeBanches[0].ShortBranchName()).To(Equal(DEVELOP_BRANCH))
		})
	})

	It("is only allowed to create from develop branch", func() {
		ForEachTestSet(features, func(feature interface{}) {
			f := feature.(AuthoredBranch)
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

	// settings like default branch
	// ----------------------------
	It("cannot be published", func() {
		ForEachTestSet(features, func(feature interface{}) {
			Expect(feature.(AuthoredBranch).CanBePublished()).To(BeFalse())
		})
	})

	It("has no publish branch", func() {
		ForEachTestSet(features, func(feature interface{}) {
			publishBranch := feature.(AuthoredBranch).PublishBranch()
			Expect(publishBranch).To(BeNil())
		})
	})

	It("has a branch name", func() {
		ForEachTestSet(features, func(feature interface{}) {
			branchName := feature.(AuthoredBranch).BranchName()
			Expect(branchName).To(Equal("refs/heads/" + FEAUTURE_BRANCH))
		})
	})

	It("has a short branch name", func() {
		ForEachTestSet(features, func(feature interface{}) {
			branchName := feature.(AuthoredBranch).ShortBranchName()
			Expect(branchName).To(Equal(FEAUTURE_BRANCH))
		})
	})
})
