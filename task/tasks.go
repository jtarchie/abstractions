package task

import (
	"sync"
	"time"
)

type Values []*value
type Tasks []*task

func (t Tasks) Await(timeout time.Duration) Values {
	var wg sync.WaitGroup

	for _, task := range t {
		wg.Add(1)
		go func() {
			defer wg.Done()
			task.Await(timeout)
		}()
	}

	wg.Wait()

	var values Values
	for _, task := range t {
		values = append(values, task.awaitWithValue(0))
	}
	return values
}
