package gitprovider_test

import (
	"log"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"

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

	It("detects the close branches and creates merge request for them", func() {
		gp.Close(mockBranch)
		// mockBranch.EXPECT().ShortBranchName().Return("mock-branch")
		// mockGitService.EXPECT().RemoteBranchExists("mock-branch").Return("", "", nil)
		// mockBranch.EXPECT().CanBeClosed().Return(true)
		// mockGitService.EXPECT().BranchList()
	})

})
