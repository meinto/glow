package cmd_test

import (
	"github.com/golang/mock/gomock"
	mockg "github.com/meinto/glow/git/mock_git"
	mockgp "github.com/meinto/glow/gitprovider/mock_gitprovider"
	. "github.com/meinto/glow/pkg/cli/cmd"
	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Root command", func() {
	var (
		mockCtrl    *gomock.Controller
		mockCommand command.Service
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockCommand = NewMockCommand(SetupRootCommand(), mockCtrl).
			SetupServices().
			Patch()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("autodetects cicd origin", func() {
		mockCommand.Cmd().SetArgs([]string{
			"--detectCicdOrigin",
		})
		mockCommand.GitProvider().(mockgp.MockServiceInterface).
			EXPECT().
			DetectCICDOrigin().
			Return("new-origin", nil)
		mockCommand.GitClient().(mockg.MockNativeServiceInterface).
			EXPECT().
			SetCICDOrigin("new-origin")
		mockCommand.Execute()
	})

	It("sets the given cicd origin", func() {
		mockCommand.Cmd().SetArgs([]string{
			"--cicdOrigin", "my-custom-origin",
		})
		mockCommand.GitClient().(mockg.MockNativeServiceInterface).
			EXPECT().
			SetCICDOrigin("my-custom-origin")
		mockCommand.Execute()
	})

})
