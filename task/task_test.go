package task_test

import (
	"errors"
	"time"

	"github.com/jtarchie/abstractions/task"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("tasks", func() {
	Context("Await", func() {
		It("returns a Task", func() {
			task := task.Async(task.NoOpFunc)
			Expect(task.Pid()).To(BeNumerically("==", 0))
		})
	})

	Context("Await", func() {
		Context("when the func fails", func() {
			It("returns the error message", func() {
				expectedError := errors.New("testing")

				task := task.Async(func() (interface{}, error) {
					return nil, expectedError
				})

				value, err := task.Await(1 * time.Millisecond)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(expectedError))
				Expect(value).To(BeNil())
			})
		})

		Context("when the func takes longer than the timeout", func() {
			It("return an error", func() {
				task := task.Async(func() (interface{}, error) {
					time.Sleep(10 * time.Second)
					return nil, nil
				})

				value, err := task.Await(1 * time.Millisecond)

				Expect(err).To(HaveOccurred())
				Expect(value).To(BeNil())
			})
		})

		Context("when the func returns a value", func() {
			It("returns that value", func() {
				task := task.Async(func() (interface{}, error) {
					return true, nil
				})

				value, err := task.Await(1 * time.Second)

				Expect(err).ToNot(HaveOccurred())
				Expect(value).To(Equal(true))
			})

			It("returns the same value when called consecutively", func() {
				task := task.Async(func() (interface{}, error) {
					return true, nil
				})

				Consistently(func() bool {
					value, err := task.Await(1 * time.Second)

					Expect(err).ToNot(HaveOccurred())
					return value.(bool)
				}).Should(Equal(true))
			})
		})
	})

	Context("Yield", func() {
		Context("when the func fails", func() {
			It("returns the error message", func() {
				expectedError := errors.New("testing")

				task := task.Async(func() (interface{}, error) {
					return nil, expectedError
				})

				value, err := task.Yield(1 * time.Millisecond)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(expectedError))
				Expect(value).To(BeNil())
			})
		})

		Context("when the func takes longer than the timeout", func() {
			It("return nothing", func() {
				task := task.Async(func() (interface{}, error) {
					time.Sleep(10 * time.Second)
					return nil, nil
				})

				value, err := task.Yield(1 * time.Millisecond)

				Expect(err).ToNot(HaveOccurred())
				Expect(value).To(BeNil())
			})
		})

		Context("when the func returns a value", func() {
			It("returns that value", func() {
				task := task.Async(func() (interface{}, error) {
					return true, nil
				})

				value, err := task.Yield(1 * time.Second)

				Expect(err).ToNot(HaveOccurred())
				Expect(value).To(Equal(true))
			})

			It("returns the same value when called consecutively", func() {
				task := task.Async(func() (interface{}, error) {
					return true, nil
				})

				Consistently(func() bool {
					value, err := task.Yield(1 * time.Second)

					Expect(err).ToNot(HaveOccurred())
					return value.(bool)
				}).Should(Equal(true))
			})
		})
	})
})
