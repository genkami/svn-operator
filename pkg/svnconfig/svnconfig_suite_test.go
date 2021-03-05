package svnconfig_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSvnconfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Svnconfig Suite")
}
