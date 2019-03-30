package semver

import (
	"testing"

	"github.com/meinto/glow/testenv"
)

func TestGetCurrentVersion(t *testing.T) {
	local, _, teardown := testenv.SetupEnv(t)
	defer teardown()

	s := setupSemverService(local.Folder)
	currentVersion, err := s.GetCurrentVersion()
	testenv.CheckForErrors(t, err)

	if currentVersion != "1.2.3" || err != nil {
		t.Errorf("current version should be '1.2.3' but is '%s'", currentVersion)
	}
}

func TestGetNextVersion(t *testing.T) {
	local, _, teardown := testenv.SetupEnv(t)
	defer teardown()

	s := setupSemverService(local.Folder)
	v, err := s.GetNextVersion("major")
	testenv.CheckForErrors(t, err)
	if v != "2.0.0" {
		t.Errorf("version should be '2.0.0' but is '%s'", v)
	}

	v, err = s.GetNextVersion("minor")
	testenv.CheckForErrors(t, err)
	if v != "1.3.0" {
		t.Errorf("version should be '1.3.0' but is '%s'", v)
	}

	v, err = s.GetNextVersion("patch")
	testenv.CheckForErrors(t, err)
	if v != "1.2.4" {
		t.Errorf("version should be '1.2.4' but is '%s'", v)
	}
}

func TestSetNextVersion(t *testing.T) {
	local, _, teardown := testenv.SetupEnv(t)
	defer teardown()

	s := setupSemverService(local.Folder)
	err := s.SetNextVersion("major")
	testenv.CheckForErrors(t, err)
	stdout, _ := local.Do("cat VERSION")
	if stdout.String() != "2.0.0" {
		t.Errorf("version should be '2.0.0' but is '%s'", stdout.String())
	}

	err = s.SetNextVersion("minor")
	testenv.CheckForErrors(t, err)
	stdout, _ = local.Do("cat VERSION")
	if stdout.String() != "2.1.0" {
		t.Errorf("version should be '2.1.0' but is '%s'", stdout.String())
	}

	err = s.SetNextVersion("patch")
	testenv.CheckForErrors(t, err)
	stdout, _ = local.Do("cat VERSION")
	if stdout.String() != "2.1.1" {
		t.Errorf("version should be '2.1.1' but is '%s'", stdout.String())
	}
}

// doesnt work yet
// func TestTagCurrentVersion(t *testing.T) {
// 	local, _, _ := testenv.SetupEnv(t)
// 	// defer teardown()

// 	s := setupSemverService(local.Folder)
// 	err := s.TagCurrentVersion()
// 	testenv.CheckForErrors(t, err)
// }

// helpers

func setupSemverService(folder string) Service {
	return NewSemverService(
		folder,
		"/bin/bash",
		"VERSION",
		"raw",
	)
}
