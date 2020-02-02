package command

import (
	"reflect"

	"github.com/meinto/glow/git"
	"github.com/meinto/glow/gitprovider"
	"github.com/meinto/glow/pkg"
	"github.com/meinto/glow/pkg/cli/cmd/internal/util"
	"github.com/spf13/cobra"
)

type Service interface {
	Cmd() *cobra.Command
	Execute() error
	Add(Service)

	SetupFlags(parent Service) Service

	SetupServices() Service
	SetGitClient(git.Service)
	GitClient() git.Service
	SetGitProvider(gitprovider.Service)
	GitProvider() gitprovider.Service

	Patch() Service
	PatchRun(fieldName string, run func(cmd Service, args []string))
}

func Setup(cmd Service) Service {
	pkg.InitGlobalConfig()
	cmd.
		SetupServices().
		Patch()

	return cmd
}

type Command struct {
	*cobra.Command
	gitClient        git.Service
	gitProvider      gitprovider.Service
	Run              func(cmd Service, args []string)
	PostRun          func(cmd Service, args []string)
	PersistentPreRun func(cmd Service, args []string)
}

func (c *Command) Cmd() *cobra.Command {
	return c.Command
}

func (c *Command) SetupFlags(parent Service) Service {
	return c
}

func (c *Command) SetGitClient(gc git.Service) {
	c.gitClient = gc
}

func (c *Command) GitClient() git.Service {
	return c.gitClient
}

func (c *Command) SetGitProvider(gp gitprovider.Service) {
	c.gitProvider = gp
}

func (c *Command) GitProvider() gitprovider.Service {
	return c.gitProvider
}

func (c *Command) SetupServices() Service {
	g, err := util.GetGitClient()
	util.ExitOnError(err)
	c.SetGitClient(g)

	gp, err := util.GetGitProvider()
	util.ExitOnError(err)
	c.SetGitProvider(gp)

	// TODO version file information from config
	// pathToRepo, _, _, err := g.GitRepoPath()
	// util.ExitOnError(err)

	// s := semver.NewSemverService(
	// 	pathToRepo,
	// 	"/bin/bash",
	// 	versionFile,
	// 	versionFileType,
	// )
	return c
}

func (c *Command) Patch() Service {
	c.PatchRun("Run", c.Run)
	c.PatchRun("PostRun", c.PostRun)
	c.PatchRun("PersistentPreRun", c.PersistentPreRun)
	return c
}

func (c *Command) PatchRun(fieldName string, run func(cmd Service, args []string)) {
	if run != nil {
		r := reflect.ValueOf(c.Command)
		f := reflect.Indirect(r).FieldByName(fieldName)

		patchedRun := func(cmd *cobra.Command, args []string) {
			run(c, args)
		}

		f.Set(reflect.ValueOf(patchedRun))
	}
}

func (c *Command) Execute() error {
	return c.Command.Execute()
}

func (c *Command) Add(cmd Service) {
	c.Command.AddCommand(cmd.Cmd())
}
