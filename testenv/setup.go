package testenv

import (
	"os"
	"testing"
)

func SetupEnv(t *testing.T) (*LocalRepository, *BareRepository, func()) {
	t.Log("init tmp git repo")

	tmpFolder := "/tmp/github.com/meinto/glow"
	local := NewLocalRepository(tmpFolder)
	bare := NewBareRepository(tmpFolder)

	bare.Init()
	local.Init(bare.Folder)

	return local, bare, func() {
		os.RemoveAll("/tmp/github.com/meinto/glow")
		t.Log("finish")
	}
}
