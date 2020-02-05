package glow_test

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/meinto/glow"
	. "github.com/meinto/glow/testutil"
)

var _ = Describe("Feature", func() {
	var branches []AuthoredBranch

	BeforeEach(func() {
		f1, _ := NewFeature("luke", "falcon")
		f2, _ := FeatureFromBranch(BRANCH_NAME_PREFIX + "feature/luke/falcon")
		f3, _ := BranchFromBranchName(BRANCH_NAME_PREFIX + "feature/luke/falcon")
		f4, _ := BranchFromBranchName("feature/luke/falcon")
		branches = []AuthoredBranch{f1, f2, f3, f4}
	})

	It("is of type feature branch", func() {
		f, _ := NewFeature("a", "b")
		rf := reflect.ValueOf(f)
		ForEachTestSet(branches, func(branch interface{}) {
			r := reflect.ValueOf(branch)
			Expect(r.Type().AssignableTo(rf.Type())).To(BeTrue())
		})
	})

	It("can be closed", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			Expect(branch.(AuthoredBranch).CanBeClosed()).To(Equal(true))
		})
	})

	It("only closes on develop", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			closeBanches := branch.(AuthoredBranch).CloseBranches([]Branch{})
			Expect(len(closeBanches)).To(Equal(1))
			Expect(closeBanches[0].ShortBranchName()).To(Equal(DEVELOP_BRANCH))
		})
	})

	It("is only allowed to create from develop branch", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			f := branch.(AuthoredBranch)
			for _, testBranch := range MockBranchCollection() {
				testBranchName := testBranch.ShortBranchName()
				if testBranchName == DEVELOP_BRANCH {
					Expect(f.CreationIsAllowedFrom(testBranch)).To(BeTrue())
				} else {
					Expect(f.CreationIsAllowedFrom(testBranch)).To(BeFalse())
				}
			}
		})
	})

	// settings like default branch
	// ----------------------------
	It("cannot be published", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			Expect(branch.(AuthoredBranch).CanBePublished()).To(BeFalse())
		})
	})

	It("has no publish branch", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			publishBranch := branch.(AuthoredBranch).PublishBranch()
			Expect(publishBranch).To(BeNil())
		})
	})

	It("has a branch name", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			branchName := branch.(AuthoredBranch).BranchName()
			Expect(branchName).To(Equal(BRANCH_NAME_PREFIX + FEAUTURE_BRANCH))
		})
	})

	It("has a short branch name", func() {
		ForEachTestSet(branches, func(branch interface{}) {
			branchName := branch.(AuthoredBranch).ShortBranchName()
			Expect(branchName).To(Equal(FEAUTURE_BRANCH))
		})
	})
})
