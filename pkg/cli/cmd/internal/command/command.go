package command

import (
	"reflect"

	"github.com/meinto/glow"
	"github.com/meinto/glow/git"
	"github.com/meinto/glow/gitprovider"
	"github.com/meinto/glow/pkg"
	"github.com/meinto/glow/pkg/cli/cmd/internal/util"
	"github.com/meinto/glow/semver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Service interface {
	Cmd() *cobra.Command
	Execute() error
	Add(Service)

	PostSetup(parent Service) Service

	SetupServices(override bool) Service
	SetGitClient(git.Service)
	GitClient() git.Service
	SetGitProvider(gitprovider.Service)
	GitProvider() gitprovider.Service
	SetSemverClient(semver.Service)
	SemverClient() semver.Service

	Patch() Service
	PatchRun(fieldName string, run func(cmd Service, args []string))

	CurrentBranch(ci bool) glow.Branch
}

func Setup(cmd Service, parent Service) Service {
	pkg.InitGlobalConfig()
	cmd.SetupServices(false)
	cmd.Patch()
	cmd.PostSetup(parent)

	return cmd
}

type Command struct {
	*cobra.Command
	gitClient         git.Service
	gitProvider       gitprovider.Service
	semverClient      semver.Service
	PreRun            func(cmd Service, args []string)
	Run               func(cmd Service, args []string)
	PostRun           func(cmd Service, args []string)
	PersistentPreRun  func(cmd Service, args []string)
	PersistentRun     func(cmd Service, args []string)
	PersistentPostRun func(cmd Service, args []string)
}

func (c *Command) Cmd() *cobra.Command {
	return c.Command
}

func (c *Command) PostSetup(parent Service) Service {
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

func (c *Command) SetSemverClient(s semver.Service) {
	c.semverClient = s
}

func (c *Command) SemverClient() semver.Service {
	return c.semverClient
}

func (c *Command) SetupServices(override bool) Service {
	if c.gitClient == nil || override {
		g, err := util.GetGitClient()
		util.ExitOnError(err)
		c.SetGitClient(g)
	}

	if c.gitProvider == nil || override {
		gp, err := util.GetGitProvider()
		util.ExitOnError(err)
		c.SetGitProvider(gp)
	}

	if c.semverClient == nil || override {
		pathToRepo, _, _, err := c.GitClient().GitRepoPath()
		util.ExitOnError(err)

		s := semver.NewSemverService(
			pathToRepo,
			"/bin/bash",
			viper.GetString("versionFile"),
			viper.GetString("versionFileType"),
		)
		c.SetSemverClient(s)
	}

	return c
}

func (c *Command) Patch() Service {
	c.PatchRun("Run", c.Run)
	c.PatchRun("PreRun", c.PreRun)
	c.PatchRun("PostRun", c.PostRun)
	c.PatchRun("PersistentPreRun", c.PersistentPreRun)
	c.PatchRun("PersistentRun", c.PersistentRun)
	c.PatchRun("PersistentPostRun", c.PersistentPostRun)
	return c
}

func (c *Command) PatchRun(fieldName string, run func(cmd Service, args []string)) {
	if run != nil {
		r := reflect.ValueOf(c.Cmd())
		f := reflect.Indirect(r).FieldByName(fieldName)

		patchedRun := func(cmd *cobra.Command, args []string) {
			run(c, args)
		}

		f.Set(reflect.ValueOf(patchedRun))
	}
}

func (c *Command) Add(cmd Service) {
	c.Command.AddCommand(cmd.Cmd())
}

func (c *Command) CurrentBranch(ci bool) glow.Branch {
	if ci {
		cb, err := c.GitProvider().GetCIBranch()
		util.ExitOnError(err)
		return cb
	} else {
		cb, _, _, err := c.GitClient().CurrentBranch()
		util.ExitOnError(err)
		return cb
	}
}
