package git_test

import (
	"strings"

	"github.com/meinto/glow/cmd"

	. "github.com/meinto/glow/git"
	"github.com/meinto/glow/testenv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("git service", func() {

	var local *testenv.LocalRepository
	var bare *testenv.BareRepository
	var teardown func()
	var s Service

	BeforeEach(func() {
		local, bare, teardown = testenv.SetupEnv()
		exec := cmd.NewCmdExecutorInDir("/bin/bash", local.Folder)
		s = NewNativeService(Options{
			CmdExecutor: exec,
		})
	})

	AfterEach(func() {
		teardown()
	})

	Describe("SetCICDOrigin", func() {
		It("sets a new origin for the local repository", func() {
			newOrigin := "https://new.origin"
			_, _, err := s.SetCICDOrigin(newOrigin)
			Expect(err).To(BeNil())
			stdout, _, err := local.Do("git remote get-url origin")
			Expect(err).To(BeNil())
			Expect(strings.TrimSpace(stdout.String())).To(Equal(newOrigin))
		})
	})

	Describe("GitRepoPath", func() {
		It("returns the path to the git repository", func() {
			exec := cmd.NewCmdExecutorInDir("/bin/bash", local.Folder+"/subfolder")
			s = NewNativeService(Options{
				CmdExecutor: exec,
			})
			repoPath, _, _, err := s.GitRepoPath()
			Expect(err).To(BeNil())
			Expect(strings.TrimPrefix(repoPath, "/private")).To(Equal(local.Folder))
		})
	})

	Describe("CurrentBranch", func() {
		It("returns 'master'", func() {
			b, _, _, err := s.CurrentBranch()
			Expect(err).To(BeNil())
			Expect(b.ShortBranchName()).To(Equal("master"))
		})

		It("returns 'test/branch'", func() {
			newBranch := "test/branch"
			local.CreateBranch(newBranch)
			local.Checkout(newBranch)
			b, _, _, err := s.CurrentBranch()
			Expect(err).To(BeNil())
			Expect(b.ShortBranchName()).To(Equal(newBranch))
		})
	})

	Describe("BranchList", func() {
		It("returns a list of all branches", func() {
			featureBranches := []string{"test/branch", "test/branch2"}
			for _, b := range featureBranches {
				local.CreateBranch(b)
			}
			bs, _, _, err := s.BranchList()
			expectedBranches := []string{"master"}
			expectedBranches = append(expectedBranches, featureBranches...)
			Expect(err).To(BeNil())
			for i, eb := range expectedBranches {
				b := bs[i]
				Expect(b.ShortBranchName()).To(Equal(eb))
			}
		})
	})

	Describe("Fetch", func() {
		It("fetches remote branches", func() {
			local2 := testenv.Clone(bare.Folder, "local2")

			local2Branch := "local2/branch"
			local2.CreateBranch(local2Branch)
			local2.Checkout(local2Branch)
			local2.Push(local2Branch)
			_, _, err := s.Fetch()
			Expect(err).To(BeNil())
			exists, branchName := local.Exists(local2Branch)
			Expect(exists).To(BeTrue())
			Expect(branchName).To(Equal(local2Branch))
		})
	})

	Describe("Stashing", func() {
		It("stashes changes", func() {
			local.Do("touch test.file")
			stdout, _, _ := local.Do("git status | grep test.file")
			Expect(strings.TrimSpace(stdout.String())).To(Equal("test.file"))
			s.AddAll()
			s.Stash()
			stdout, _, _ = local.Do("git status | grep test.file")
			Expect(strings.TrimSpace(stdout.String())).To(BeEmpty())
		})

		It("pops the stash", func() {
			local.Do("touch test.file")
			stdout, _, _ := local.Do("git status | grep test.file")
			Expect(strings.TrimSpace(stdout.String())).To(Equal("test.file"))
			s.AddAll()
			s.Stash()
			stdout, _, _ = local.Do("git status | grep test.file")
			Expect(strings.TrimSpace(stdout.String())).To(BeEmpty())
			s.StashPop()
			stdout, _, _ = local.Do("git status | grep test.file")
			Expect(strings.TrimSpace(stdout.String())).NotTo(BeEmpty())
		})
	})

	Describe("Commit", func() {
		It("commits the changes", func() {
			local.Do("touch test.file")
			s.AddAll()
			_, _, err := s.Commit("Commit test.file")
			Expect(err).To(BeNil())
			stdout, _, _ := local.Do(`git rev-list --left-only --count master...origin/master`)
			Expect(strings.TrimSpace(stdout.String())).To(Equal("1"))
		})
	})

	Describe("Push", func() {
		It("pushes changes", func() {
			local.Do("touch test.file")
			s.AddAll()
			s.Commit("Commit test.file")
			stdout, _, _ := local.Do(`git rev-list --left-only --count master...origin/master`)
			Expect(strings.TrimSpace(stdout.String())).To(Equal("1"))
			_, _, err := s.Push(false)
			Expect(err).To(BeNil())
			stdout, _, _ = local.Do(`git rev-list --left-only --count master...origin/master`)
			Expect(strings.TrimSpace(stdout.String())).To(Equal("0"))
		})
	})
})
