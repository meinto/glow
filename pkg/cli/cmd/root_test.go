package cmd_test

import (
	"github.com/golang/mock/gomock"
	mockg "github.com/meinto/glow/git/__mock__"
	mockgp "github.com/meinto/glow/gitprovider/__mock__"
	. "github.com/meinto/glow/pkg/cli/cmd"
	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Root command", func() {
	var (
		mockCtrl        *gomock.Controller
		mockCommand     command.Service
		mockGitProvider mockgp.MockServiceInterface
		mockGitClient   mockg.MockNativeServiceInterface
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockCommand = NewMockCommand(SetupRootCommand(), mockCtrl)
		mockGitProvider = mockCommand.GitProvider().(mockgp.MockServiceInterface)
		mockGitClient = mockCommand.GitClient().(mockg.MockNativeServiceInterface)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("autodetects cicd origin", func() {
		mockCommand.Cmd().SetArgs([]string{
			"--detectCicdOrigin",
		})
		mockGitProvider.EXPECT().DetectCICDOrigin().Return("new-origin", nil)
		mockGitClient.EXPECT().SetCICDOrigin("new-origin")
		mockCommand.Execute()
	})

	It("sets the given cicd origin", func() {
		mockCommand.Cmd().SetArgs([]string{
			"--cicdOrigin", "my-custom-origin",
		})
		mockGitClient.EXPECT().SetCICDOrigin("my-custom-origin")
		mockCommand.Execute()
	})

})
