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
		mockCtrl     *gomock.Controller
		mockRootCmd  command.Service
		mockCloseCmd command.Service
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRootCmd = NewMockCommand(SetupRootCommand(), mockCtrl).
			SetupServices(true).
			Patch()
		mockCloseCmd = NewMockCommand(SetupCloseCommand(mockRootCmd), mockCtrl).
			SetupServices(true).
			Patch()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("trys to close the current branch", func() {
		mockRootCmd.Cmd().SetArgs([]string{
			"close",
		})
		b := glow.NewBranch("test")
		mockCloseCmd.GitClient().(mockg.MockNativeServiceInterface).
			EXPECT().
			CurrentBranch().
			Return(b, "", "", nil)
		mockCloseCmd.GitProvider().(mockgp.MockServiceInterface).
			EXPECT().
			Close(b)
		mockRootCmd.Execute()
	})

	It("trys to detect the cicd branch and close it", func() {
		mockRootCmd.Cmd().SetArgs([]string{
			"close", "--ci",
		})
		b := glow.NewBranch("test")
		mockCloseCmd.GitProvider().(mockgp.MockServiceInterface).
			EXPECT().
			GetCIBranch().
			Return(b, nil)
		mockCloseCmd.GitProvider().(mockgp.MockServiceInterface).
			EXPECT().
			Close(b)
		mockRootCmd.Execute()
	})

})
