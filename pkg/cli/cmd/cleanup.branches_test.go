package cmd_test

import (
	"github.com/golang/mock/gomock"
	mockg "github.com/meinto/glow/git/__mock__"
	. "github.com/meinto/glow/pkg/cli/cmd"
	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Cleanup branches command", func() {
	var (
		mockCtrl               *gomock.Controller
		mockRootCmd            command.Service
		mockCleanupCmd         command.Service
		mockCleanupBranchesCmd command.Service
		mockGitClient          mockg.MockNativeServiceInterface
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRootCmd = NewMockCommand(SetupRootCommand(), mockCtrl).
			SetupServices(true).
			Patch()
		mockCleanupCmd = NewMockCommand(SetupCleanupCommand(mockRootCmd), mockCtrl).
			SetupServices(true).
			Patch()
		mockCleanupBranchesCmd = NewMockCommand(SetupCleanupBranchesCommand(mockCleanupCmd), mockCtrl).
			SetupServices(true).
			Patch()
		mockGitClient = mockCleanupBranchesCmd.GitClient().(mockg.MockNativeServiceInterface)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("cleans up the branches #default", func() {
		mockRootCmd.Cmd().SetArgs([]string{
			"cleanup", "branches",
		})
		mockGitClient.EXPECT().CleanupBranches(false, false)
		mockRootCmd.Execute()
	})

	It("cleans up gone branches", func() {
		mockRootCmd.Cmd().SetArgs([]string{
			"cleanup", "branches", "--gone",
		})
		mockGitClient.EXPECT().CleanupBranches(true, false)
		mockRootCmd.Execute()
	})

	It("cleans up untracked branches", func() {
		mockRootCmd.Cmd().SetArgs([]string{
			"cleanup", "branches", "--untracked",
		})
		mockGitClient.EXPECT().CleanupBranches(false, true)
		mockRootCmd.Execute()
	})

})
