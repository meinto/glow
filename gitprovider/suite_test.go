package gitprovider_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGlow(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Git Provider Suite")
}
