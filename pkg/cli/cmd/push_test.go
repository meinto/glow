package cmd_test

import (
	"github.com/golang/mock/gomock"
	"github.com/meinto/glow"
	mockg "github.com/meinto/glow/git/__mock__"

	// mockgp "github.com/meinto/glow/gitprovider/__mock__"
	. "github.com/meinto/glow/pkg/cli/cmd"
	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Push command", func() {
	var (
		mockCtrl    *gomock.Controller
		mockRootCmd command.Service
		mockPushCmd command.Service
		// mockGitProvider mockgp.MockServiceInterface
		mockGitClient mockg.MockNativeServiceInterface
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRootCmd = NewMockCommand(SetupRootCommand(), mockCtrl)
		mockPushCmd = NewMockCommand(SetupPushCommand(mockRootCmd), mockCtrl)
		// mockGitProvider = mockPushCmd.GitProvider().(mockgp.MockServiceInterface)
		mockGitClient = mockPushCmd.GitClient().(mockg.MockNativeServiceInterface)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("pushes the branch #happyPath", func() {
		mockRootCmd.Cmd().SetArgs([]string{
			"push",
		})
		b := glow.NewBranch("test")
		mockGitClient.EXPECT().CurrentBranch().Return(b, "", "", nil)
		mockGitClient.EXPECT().RemoteBranchExists(b.BranchName()).Return(true, "", "", nil)
		setUpstream := false
		mockGitClient.EXPECT().Push(setUpstream)
		mockRootCmd.Execute()
	})

	It("pushes the branch and sets the upstream", func() {
		mockRootCmd.Cmd().SetArgs([]string{
			"push",
		})
		b := glow.NewBranch("test")
		mockGitClient.EXPECT().CurrentBranch().Return(b, "", "", nil)
		mockGitClient.EXPECT().RemoteBranchExists(b.BranchName()).Return(false, "", "", nil)
		setUpstream := true
		mockGitClient.EXPECT().Push(setUpstream)
		mockRootCmd.Execute()
	})

})
