package semver

import (
	"testing"

	"github.com/meinto/glow/testenv"
)

func setupSemverService(folder string) Service {
	return NewSemverService(
		folder,
		"/usr/local/bin/git",
		"VERSION",
		"raw",
	)
}
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
	if v != "2.0.0" || err != nil {
		t.Errorf("version should be '2.0.0' but is '%s'", v)
	}

	v, err = s.GetNextVersion("minor")
	testenv.CheckForErrors(t, err)
	if v != "1.3.0" || err != nil {
		t.Errorf("version should be '1.3.0' but is '%s'", v)
	}

	v, err = s.GetNextVersion("patch")
	testenv.CheckForErrors(t, err)
	if v != "1.2.4" || err != nil {
		t.Errorf("version should be '1.2.4' but is '%s'", v)
	}
}
