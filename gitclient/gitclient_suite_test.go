package gitclient_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGitclient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gitclient Suite")
}
