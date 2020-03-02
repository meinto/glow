package cmd_test

import (
	"github.com/golang/mock/gomock"
	"github.com/meinto/glow"
	mockg "github.com/meinto/glow/git/__mock__"
	mockgp "github.com/meinto/glow/gitprovider/__mock__"
	. "github.com/meinto/glow/pkg/cli/cmd"
	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Cleanup branches command", func() {
	var (
		mockCtrl        *gomock.Controller
		mockRootCmd     command.Service
		mockCloseCmd    command.Service
		mockGitProvider mockgp.MockServiceInterface
		mockGitClient   mockg.MockNativeServiceInterface
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRootCmd = NewMockCommand(SetupRootCommand(), mockCtrl)
		mockCloseCmd = NewMockCommand(SetupCloseCommand(mockRootCmd), mockCtrl)
		mockGitProvider = mockCloseCmd.GitProvider().(mockgp.MockServiceInterface)
		mockGitClient = mockCloseCmd.GitClient().(mockg.MockNativeServiceInterface)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("trys to close the current branch", func() {
		mockRootCmd.Cmd().SetArgs([]string{
			"close",
		})
		b := glow.NewBranch("test")
		mockGitClient.EXPECT().CurrentBranch().Return(b, "", "", nil)
		mockGitProvider.EXPECT().Close(b)
		mockRootCmd.Execute()
	})

	It("trys to detect the cicd branch and close it", func() {
		mockRootCmd.Cmd().SetArgs([]string{
			"close", "--ci",
		})
		b := glow.NewBranch("test")
		mockGitProvider.EXPECT().GetCIBranch().Return(b, nil)
		mockGitProvider.EXPECT().Close(b)
		mockRootCmd.Execute()
	})

})
