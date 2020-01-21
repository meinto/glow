package testenv

import (
	"os"
)

func SetupEnv() (*LocalRepository, *BareRepository, func()) {
	tmpFolder := "/tmp/github.com/meinto/glow"
	local := NewLocalRepository(tmpFolder)
	bare := NewBareRepository(tmpFolder)

	bare.Init()
	local.Init(bare.Folder)

	return local, bare, func() {
		os.RemoveAll("/tmp/github.com/meinto/glow")
	}
}
