package cmd_test

import (
	"log"

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
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRootCommand = NewMockCommand(SetupRootCommand(), mockCtrl).
			SetupServices(true).
			Patch()
		releaseMockCommand = NewMockCommand(SetupReleaseCommand(mockRootCommand), mockCtrl).
			SetupServices(true).
			Patch()
		log.Println(releaseMockCommand)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("creates a release branch with the current version", func() {
		mockRootCommand.Cmd().SetArgs([]string{
			"release", "current",
		})
		CURRENT_VERSION := "2.2.2"
		releaseMockCommand.SemverClient().(mocksemver.MockServiceInterface).
			EXPECT().
			GetCurrentVersion().
			Return(CURRENT_VERSION, nil)
		b, _ := glow.NewRelease(CURRENT_VERSION)
		releaseMockCommand.GitClient().(mockg.MockNativeServiceInterface).
			EXPECT().
			Create(b, false)
		releaseMockCommand.GitClient().(mockg.MockNativeServiceInterface).
			EXPECT().
			Checkout(b)
		releaseMockCommand.SemverClient().(mocksemver.MockServiceInterface).
			EXPECT().
			SetVersion(CURRENT_VERSION)
		mockRootCommand.Execute()
	})
})
