package task

import (
	"errors"
	"sync"
	"time"
)

type work func() (interface{}, error)
type retval struct {
	err   error
	value interface{}
}
type task struct {
	fn    work
	wait  chan retval
	final *retval
	once  sync.Once
}

var NoOpFunc = func() (interface{}, error) {
	return nil, nil
}

func Async(fn work) *task {
	task := &task{
		fn:   fn,
		wait: make(chan retval, 1),
	}

	go func() {
		value, err := task.fn()
		task.wait <- retval{
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
	t.once.Do(func() {
		select {
		case retval := <-t.wait:
			t.final = &retval
		case <-time.After(timeout):
			t.final = &retval{
				value: nil,
				err:   errors.New("timeout occurred"),
			}
		}
	})

	return t.final.value, t.final.err
}

func (t *task) Yield(timeout time.Duration) (interface{}, error) {
	if t.final == nil {
		select {
		case retval := <-t.wait:
			t.final = &retval
		case <-time.After(timeout):
			return nil, nil
		}
	}

	return t.final.value, t.final.err
}
