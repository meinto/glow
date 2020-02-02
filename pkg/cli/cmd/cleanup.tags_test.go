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
		mockCtrl           *gomock.Controller
		mockRootCmd        command.Service
		mockCleanupCmd     command.Service
		mockCleanupTagsCmd command.Service
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRootCmd = NewMockCommand(SetupRootCommand(), mockCtrl).
			SetupServices().
			Patch()
		mockCleanupCmd = NewMockCommand(SetupCleanupCommand(mockRootCmd), mockCtrl).
			SetupServices().
			Patch()
		mockCleanupTagsCmd = NewMockCommand(SetupCleanupTagsCommand(mockCleanupCmd), mockCtrl).
			SetupServices().
			Patch()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("cleans up the tags #default", func() {
		mockRootCmd.Cmd().SetArgs([]string{
			"cleanup", "tags",
		})
		mockCleanupTagsCmd.GitClient().(mockg.MockNativeServiceInterface).
			EXPECT().
			CleanupTags(false)
		mockRootCmd.Execute()
	})

	It("cleans up untracked tags", func() {
		mockRootCmd.Cmd().SetArgs([]string{
			"cleanup", "tags", "--untracked",
		})
		mockCleanupTagsCmd.GitClient().(mockg.MockNativeServiceInterface).
			EXPECT().
			CleanupTags(true)
		mockRootCmd.Execute()
	})

})
