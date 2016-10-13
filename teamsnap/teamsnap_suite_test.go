package teamsnap_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTeamsnap(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Teamsnap Suite")
}
