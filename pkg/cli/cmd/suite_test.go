package cmd_test

import (
	"testing"

	// . "github.com/meinto/glow"
	"github.com/golang/mock/gomock"
	"github.com/meinto/glow/git"
	mockg "github.com/meinto/glow/git/__mock__"
	mockgp "github.com/meinto/glow/gitprovider/__mock__"
	"github.com/meinto/glow/pkg/cli/cmd/internal/command"
	mocksemver "github.com/meinto/glow/semver/__mock__"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGlow(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Package CLI Cmd Suite")
}

type MockGitClient struct {
	git.Service
}

type MockCommand struct {
	command.Service
	mockCtrl *gomock.Controller
}

func (tc *MockCommand) SetupServices(override bool) command.Service {
	tc.SetGitClient(mockg.NewMockNativeService(tc.mockCtrl))
	tc.SetGitProvider(mockgp.NewMockService(tc.mockCtrl))
	tc.SetSemverClient(mocksemver.NewMockService(tc.mockCtrl))
	return tc
}

func NewMockCommand(originalCommand command.Service, mockCtrl *gomock.Controller) *MockCommand {
	cmd := &MockCommand{originalCommand, mockCtrl}
	cmd.SetupServices(true)
	cmd.Patch()
	return cmd
}
