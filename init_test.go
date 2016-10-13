package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

const serverAddress = "localhost:8000"

func TestTeam(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Team Suite")
}
