package main

//go:generate ../../.circleci/generate-assets.sh

import (
	"github.com/meinto/glow/pkg"
	clicmd "github.com/meinto/glow/pkg/cli/cmd"
)

func main() {
	pkg.InitGlobalConfig()
	clicmd.Execute()
}
