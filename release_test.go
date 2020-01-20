package glow_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/meinto/glow"
	. "github.com/meinto/glow/testutil"
)

var _ = Describe("Release", func() {
	var features []Release

	BeforeEach(func() {
		f1, _ := NewRelease("1.2.4")
		f2, _ := ReleaseFromBranch("release/v1.2.4")
		features = []Release{f1, f2}
	})

	It("can be closed", func() {
		ForEachTestSet(features, func(feature interface{}) {
			Expect(feature.(Release).CanBeClosed()).To(Equal(true))
		})
	})

	It("can be published", func() {
		ForEachTestSet(features, func(feature interface{}) {
			Expect(feature.(Release).CanBePublished()).To(Equal(true))
		})
	})

	It("only closes on release branch", func() {
		ForEachTestSet(features, func(feature interface{}) {
			closeBanches := feature.(Release).CloseBranches(MockBranchCollection())
			Expect(len(closeBanches)).To(Equal(1))
			Expect(closeBanches[0].ShortBranchName()).To(Equal(DEVELOP_BRANCH))
		})
	})

	It("is only allowed to create from develop branch", func() {
		ForEachTestSet(features, func(feature interface{}) {
			f := feature.(Release)
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
		ForEachTestSet(features, func(feature interface{}) {
			f := feature.(Release)
			Expect(f.PublishBranch().ShortBranchName()).To(Equal(MASTER_BRANCH))
		})
	})
})
