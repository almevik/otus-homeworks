package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(tasks <-chan Task, errors chan error) {
	for task := range tasks {
		err := task()

		if err != nil {
			select {
			case errors <- err:
			default:
				return
			}
		}
	}
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
func Run(tasks []Task, N int, M int) error {
	if N <= 0 || len(tasks) == 0 {
		return nil
	}

	wg := sync.WaitGroup{}
	wg.Add(N)

	tCh := make(chan Task, len(tasks))

	for _, task := range tasks {
		tCh <- task
	}

	close(tCh)

	errCh := make(chan error, M)

	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			worker(tCh, errCh)
		}()
	}

	wg.Wait()
	close(errCh)

	if len(errCh) == M {
		return ErrErrorsLimitExceeded
	}

	return nil
}

