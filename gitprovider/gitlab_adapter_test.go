package gitprovider_test

import (
	"log"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"

	"github.com/meinto/glow"
	mockb "github.com/meinto/glow/__mock__"
	mockg "github.com/meinto/glow/git/__mock__"
	"github.com/meinto/glow/gitprovider"
	. "github.com/meinto/glow/gitprovider"
)

var _ = Describe("Branch", func() {
	var gp Service
	var mockBranch mockb.MockBranchInterface
	var mockGitService mockg.MockNativeServiceInterface

	BeforeEach(func() {
		mockCtrl := gomock.NewController(GinkgoT())
		newMockGitService := mockg.NewMockNativeService(mockCtrl)
		mockBranch = mockb.NewMockBranch(mockCtrl)
		gp = NewGitlabService(gitprovider.Options{
			"mock-endpoint",
			"mock-namespace",
			"mock-project",
			"mock-token",
			false,
		})
		gp.SetGitService(newMockGitService)
		mockGitService = gp.GitService().(mockg.MockNativeServiceInterface)
		log.Println(mockBranch, mockGitService)
	})

	Context("close command", func() {
		TestCreateMergeRequest := func() {
			mockBranch.EXPECT().ShortBranchName().Return("source-branch")
			// TODO: Record http request
		}

		It("detects the close branches and creates merge request for them", func() {
			mockBranch.EXPECT().ShortBranchName().Return("mock-branch")
			mockGitService.EXPECT().RemoteBranchExists("mock-branch").Return("", "", nil)
			mockBranch.EXPECT().CanBeClosed().Return(true)
			branchList := []glow.Branch{
				glow.NewBranch("target-1"),
				glow.NewBranch("target-2"),
				glow.NewBranch("branch-3"),
			}
			mockGitService.EXPECT().BranchList().Return(branchList, "", "", nil)
			targetBranches := []glow.Branch{
				glow.NewBranch("target-1"),
				glow.NewBranch("target-2"),
			}
			mockBranch.EXPECT().CloseBranches(branchList).Return(targetBranches)
			TestCreateMergeRequest()
			gp.Close(mockBranch)
		})
	})

})
