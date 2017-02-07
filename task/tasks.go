package task

import (
	"sync"
	"time"
)

type Values []*value
type Tasks []*task

func (tasks Tasks) Await(timeout time.Duration) Values {
	var wg sync.WaitGroup

	for _, t := range tasks {
		wg.Add(1)
		go func(t *task) {
			defer wg.Done()
			t.Await(timeout)
		}(t)
	}

	wg.Wait()

	var values Values
	for _, t := range tasks {
		values = append(values, t.awaitWithValue(0))
	}
	return values
}
