package task_test

import (
	"time"

	"github.com/jtarchie/abstractions/task"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tasks", func() {
	Context("Await", func() {
		Context("when one func returns and another timesout", func() {
			It("returns the appropriate value and timeout error", func() {
				tasks := task.Tasks{
					task.Async(func() (interface{}, error) {
						return true, nil
					}),
					task.Async(func() (interface{}, error) {
						time.Sleep(10 * time.Second)
						return nil, nil
					}),
				}

				values := tasks.Await(1 * time.Millisecond)

				Expect(values[0].Returned).To(Equal(true))
				Expect(values[0].Err).ToNot(HaveOccurred())

				Expect(values[1].Returned).To(BeNil())
				Expect(values[1].Err).To(HaveOccurred())
			})
		})
	})
})
