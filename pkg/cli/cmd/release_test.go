package cmd_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/meinto/glow/pkg/cli/cmd"
	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Release command", func() {
	var (
		mockCtrl           *gomock.Controller
		mockRootCommand    command.Service
		releaseMockCommand command.Service
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRootCommand = NewMockCommand(CreateRootCmd(), mockCtrl)
		mockRootCommand.Init().Patch()
		releaseMockCommand = NewMockCommand(ReleaseCmd(mockRootCommand), mockCtrl)
		releaseMockCommand.Init().Patch()
		mockRootCommand.Add(releaseMockCommand)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	// It("pushes the changes", func() {
	// 	mockRootCommand.Cmd().SetArgs([]string{
	// 		"release", "2.3.4",
	// 	})
	// 	releaseMockCommand.GitClient().(mockg.MockNativeServiceInterface).
	// 		EXPECT().
	// 		GitRepoPath().
	// 		Return("path-to-repo", "", "", nil)
	// 	releaseMockCommand.GitClient().(mockg.MockNativeServiceInterface).
	// 		EXPECT().
	// 		Create(nil, false)
	// 	mockRootCommand.Execute()
	// 	// mockCommand.Cmd().Run(mockCommand.Cmd(), []string{
	// 	// 	"release", "2.3.4",
	// 	// })
	// })
})
