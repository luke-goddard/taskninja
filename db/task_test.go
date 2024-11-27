package db

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestTasksTable(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Task Database Suite")
}

var _ = Describe("Task Pretty Age", func() {
	DescribeTable("Pretty Age", func(time time.Duration, expect string) {
		task := Task{}
		result := task.PrettyAge(time)
		Expect(result).To(Equal(expect))
	},
		Entry("0m", time.Duration(0)*time.Second, "0m"),
		Entry("1m", time.Duration(1)*time.Second, "0m"),
		Entry("1m", time.Duration(60)*time.Second, "1m"),
		Entry("1m1s", time.Duration(61)*time.Second, "1m"),
		Entry("1h0m", time.Duration(60)*time.Minute, "1h0m"),
		Entry("1h1m", time.Duration(61)*time.Minute, "1h1m"),
		Entry("1d0h", time.Duration(24)*time.Hour, "1d0h"),
		Entry("1d1h", time.Duration(25)*time.Hour, "1d1h"),
		Entry("1d4h", time.Duration(28)*time.Hour, "1d4h"),
		Entry("1w0d", time.Duration(24)*time.Hour*7, "1w0d"),
		Entry("3w3d", time.Duration(1)*time.Hour*526, "3w3d"),
	)
})
