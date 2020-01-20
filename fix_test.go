package glow_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/meinto/glow"
	. "github.com/meinto/glow/testutil"
)

var _ = Describe("Fix", func() {
	var features []Fix

	BeforeEach(func() {
		f1, _ := NewFix("luke", "falcon")
		f2, _ := FixFromBranch("fix/luke/falcon")
		features = []Fix{f1, f2}
	})

	It("can be closed", func() {
		ForEachTestSet(features, func(feature interface{}) {
			Expect(feature.(Fix).CanBeClosed()).To(Equal(true))
		})
	})

	It("only closes on release branch", func() {
		ForEachTestSet(features, func(feature interface{}) {
			closeBanches := feature.(Fix).CloseBranches(MockBranchCollection())
			Expect(len(closeBanches)).To(Equal(1))
			Expect(closeBanches[0].ShortBranchName()).To(Equal(RELEASE_BRANCH))
		})
	})

	It("is only allowed to create from release branch", func() {
		ForEachTestSet(features, func(feature interface{}) {
			f := feature.(Fix)
			for _, testBranch := range MockBranchCollection() {
				testBranchName := testBranch.ShortBranchName()
				if testBranchName == RELEASE_BRANCH {
					Expect(f.CreationIsAllowedFrom(testBranchName)).To(BeTrue())
				} else {
					Expect(f.CreationIsAllowedFrom(testBranchName)).To(BeFalse())
				}
			}
		})
	})
})
