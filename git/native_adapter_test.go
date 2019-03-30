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

func TestGitRepoPath(t *testing.T) {
	local, _, teardown := testenv.SetupEnv(t)
	defer teardown()

	s := setupNativeGitService(local.Folder + "/subfolder")
	repoPath, err := s.GitRepoPath()
	testenv.CheckForErrors(t, err)

	if strings.TrimPrefix(repoPath, "/private") != local.Folder {
		t.Errorf("repo path should be %s but is %s", local.Folder, repoPath)
	}
}

func TestCurrentBranch(t *testing.T) {
	local, _, teardown := testenv.SetupEnv(t)
	defer teardown()

	s := setupNativeGitService(local.Folder)
	b, err := s.CurrentBranch()
	testenv.CheckForErrors(t, err)
	if b.ShortBranchName() != "master" {
		t.Errorf("branch should be 'master' but is '%s'", b.ShortBranchName())
	}

	newBranch := "test/branch"
	local.CreateBranch(newBranch)
	local.Checkout(newBranch)

	b, err = s.CurrentBranch()
	testenv.CheckForErrors(t, err)
	if b.ShortBranchName() != newBranch {
		t.Errorf("branch should be 'master' but is '%s'", b.ShortBranchName())
	}
}
