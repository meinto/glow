package git

import (
	"strings"
	"testing"

	"github.com/meinto/glow/cmd"

	"github.com/meinto/glow/testenv"
)

func setupNativeGitService(pathToRepo string) Service {
	exec := cmd.NewCmdExecutorInDir("/bin/bash", pathToRepo)
	s := NewNativeService(exec)
	return s
}

func TestSetCICDOrigin(t *testing.T) {
	local, _, teardown := testenv.SetupEnv(t)
	defer teardown()

	newOrigin := "https://new.origin"
	s := setupNativeGitService(local.Folder)
	err := s.SetCICDOrigin(newOrigin)
	testenv.CheckForErrors(t, err)

	stdout, err := local.Do("git remote get-url origin")
	testenv.CheckForErrors(t, err)
	if strings.TrimSpace(stdout.String()) != newOrigin {
		t.Errorf("origin should be %s but is %s", newOrigin, stdout.String())
	}
}
