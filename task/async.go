package task

import (
	"errors"
	"time"
)

type work func() (interface{}, error)
type retval struct {
	err   error
	value interface{}
}
type task struct {
	fn     work
	retval chan retval
}

var NoOpFunc = func() (interface{}, error) {
	return nil, nil
}

func Async(fn work) *task {
	task := &task{
		fn:     fn,
		retval: make(chan retval, 1),
	}

	go func() {
		value, err := task.fn()
		task.retval <- retval{
			err:   err,
			value: value,
		}
	}()

	return task
}

func (t *task) Pid() uint64 {
	return 0
}

func (t *task) Await(timeout time.Duration) (interface{}, error) {
	select {
	case retval := <-t.retval:
		return retval.value, retval.err
	case <-time.After(timeout):
		return nil, errors.New("timeout occurred")
	}
}
