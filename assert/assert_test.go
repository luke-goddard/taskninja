package assert_test

import (
	"os"
	"testing"

	"github.com/luke-goddard/taskninja/assert"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestAsserts(t *testing.T) {
	log.Logger = log.Output(zerolog.Nop())
	RegisterFailHandler(Fail)
	RunSpecs(t, "Asserts Suite")
}

var _ = Describe("Assert", func() {
	BeforeEach(func() { os.Setenv(assert.TaskNinjaSkipAssert, "true") })
	AfterEach(func() { os.Setenv(assert.TaskNinjaSkipAssert, "") })
	It("should assert", func() {
		assert.True(true, "This is true")
		assert.True(false, "This is false")
	})
})
