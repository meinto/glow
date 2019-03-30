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
	currentVersion, _ := s.GetCurrentVersion()

	if currentVersion != "1.2.3" {
		t.Errorf("current version should be '1.2.3' but is '%s'", currentVersion)
	}
}
