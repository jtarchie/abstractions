package task

import (
	"errors"
	"sync"
	"time"
)

type work func() (interface{}, error)
type value struct {
	Err      error
	Returned interface{}
}
type task struct {
	fn    work
	wait  chan value
	final *value
	once  sync.Once
}

var NoOpFunc = func() (interface{}, error) {
	return nil, nil
}

func Async(fn work) *task {
	task := &task{
		fn:   fn,
		wait: make(chan value, 1),
	}

	go func() {
		returned, err := task.fn()
		task.wait <- value{
			Err:      err,
			Returned: returned,
		}
	}()

	return task
}

func (t *task) Pid() uint64 {
	return 0
}

func (t *task) awaitWithValue(timeout time.Duration) *value {
	t.once.Do(func() {
		returned, err := t.Yield(timeout)
		if returned == nil && err == nil {
			t.final = &value{
				Returned: nil,
				Err:      errors.New("timeout occurred"),
			}
		}
	})
	return t.final
}

func (t *task) Await(timeout time.Duration) (interface{}, error) {
	t.awaitWithValue(timeout)
	return t.final.Returned, t.final.Err
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

	return t.final.Returned, t.final.Err
}
