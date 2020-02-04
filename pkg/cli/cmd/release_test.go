package cmd_test

import (
	"github.com/golang/mock/gomock"
	"github.com/meinto/glow"
	mockg "github.com/meinto/glow/git/__mock__"
	. "github.com/meinto/glow/pkg/cli/cmd"
	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
	mocksemver "github.com/meinto/glow/semver/__mock__"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Release command", func() {
	var (
		mockCtrl           *gomock.Controller
		mockRootCommand    command.Service
		releaseMockCommand command.Service
		mockGitClient      mockg.MockNativeServiceInterface
		mockSemverClient   mocksemver.MockServiceInterface
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRootCommand = NewMockCommand(SetupRootCommand(), mockCtrl).
			SetupServices(true).
			Patch()
		releaseMockCommand = NewMockCommand(SetupReleaseCommand(mockRootCommand), mockCtrl).
			SetupServices(true).
			Patch()
		mockGitClient = releaseMockCommand.GitClient().(mockg.MockNativeServiceInterface)
		mockSemverClient = releaseMockCommand.SemverClient().(mocksemver.MockServiceInterface)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("creates a release branch with the current version", func() {
		mockRootCommand.Cmd().SetArgs([]string{
			"release", "current",
		})
		CURRENT_VERSION := "2.2.2"
		mockSemverClient.EXPECT().GetCurrentVersion().Return(CURRENT_VERSION, nil)
		b, _ := glow.NewRelease(CURRENT_VERSION)
		mockGitClient.EXPECT().Create(b, false)
		mockGitClient.EXPECT().Checkout(b)
		mockSemverClient.EXPECT().SetVersion(CURRENT_VERSION)
		mockRootCommand.Execute()
	})
})
