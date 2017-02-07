package abstractions_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAbstractions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Abstractions Suite")
}
