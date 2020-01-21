package semver_test

import (
	"strings"

	. "github.com/meinto/glow/semver"
	"github.com/meinto/glow/testenv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("semver service", func() {

	var local *testenv.LocalRepository
	var bare *testenv.BareRepository
	var teardown func()
	var s Service

	BeforeEach(func() {
		local, bare, teardown = testenv.SetupEnv()
		s = setupSemverService(local.Folder)
		s = NewLoggingService(s)
	})

	AfterEach(func() {
		teardown()
	})

	Describe("GetCurrentVersion", func() {
		It("returns the current version", func() {
			currentVersion, err := s.GetCurrentVersion()
			Expect(err).To(BeNil())
			Expect(currentVersion).To(Equal("1.2.3"))
		})
	})

	Describe("GetNextVersion", func() {
		It("returns the next patch version", func() {
			version, err := s.GetNextVersion("patch")
			Expect(err).To(BeNil())
			Expect(version).To(Equal("1.2.4"))
		})

		It("returns the next minor version", func() {
			version, err := s.GetNextVersion("minor")
			Expect(err).To(BeNil())
			Expect(version).To(Equal("1.3.0"))
		})

		It("returns the next major version", func() {
			version, err := s.GetNextVersion("major")
			Expect(err).To(BeNil())
			Expect(version).To(Equal("2.0.0"))
		})
	})

	Describe("SetNextVersion", func() {
		It("sets the next patch version", func() {
			err := s.SetNextVersion("patch")
			stdout, _, _ := local.Do("cat VERSION")
			Expect(err).To(BeNil())
			Expect(stdout.String()).To(Equal("1.2.4"))
		})

		It("sets the next minor version", func() {
			err := s.SetNextVersion("minor")
			stdout, _, _ := local.Do("cat VERSION")
			Expect(err).To(BeNil())
			Expect(stdout.String()).To(Equal("1.3.0"))
		})

		It("sets the next major version", func() {
			err := s.SetNextVersion("major")
			stdout, _, _ := local.Do("cat VERSION")
			Expect(err).To(BeNil())
			Expect(stdout.String()).To(Equal("2.0.0"))
		})
	})

	Describe("TagCurrentVersion", func() {
		It("sets a git tag for the current version", func() {
			err := s.TagCurrentVersion()
			stdout, _, _ := bare.Do("git tag | grep v1.2.3")
			Expect(err).To(BeNil())
			Expect(strings.TrimSpace(stdout.String())).To(Equal("v1.2.3"))
		})
	})
})

// helpers
func setupSemverService(folder string) Service {
	return NewSemverService(
		folder,
		"/bin/bash",
		"VERSION",
		"raw",
	)
}
