package cmd_test

import (
	"github.com/golang/mock/gomock"
	mockg "github.com/meinto/glow/git/__mock__"
	. "github.com/meinto/glow/pkg/cli/cmd"
	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Cleanup branches command", func() {
	var (
		mockCtrl           *gomock.Controller
		mockRootCmd        command.Service
		mockCleanupCmd     command.Service
		mockCleanupTagsCmd command.Service
		mockGitClient      mockg.MockNativeServiceInterface
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRootCmd = NewMockCommand(SetupRootCommand(), mockCtrl)
		mockCleanupCmd = NewMockCommand(SetupCleanupCommand(mockRootCmd), mockCtrl)
		mockCleanupTagsCmd = NewMockCommand(SetupCleanupTagsCommand(mockCleanupCmd), mockCtrl)
		mockGitClient = mockCleanupTagsCmd.GitClient().(mockg.MockNativeServiceInterface)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("cleans up the tags #default", func() {
		mockRootCmd.Cmd().SetArgs([]string{
			"cleanup", "tags",
		})
		mockGitClient.EXPECT().CleanupTags(false)
		mockRootCmd.Execute()
	})

	It("cleans up untracked tags", func() {
		mockRootCmd.Cmd().SetArgs([]string{
			"cleanup", "tags", "--untracked",
		})
		mockGitClient.EXPECT().CleanupTags(true)
		mockRootCmd.Execute()
	})

})
