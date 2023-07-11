package workers

import "sync"

type Workers[T any] struct {
	Work    chan func() T
	Results chan T
	wg      sync.WaitGroup
}

func New[T any](numWorkers int) *Workers[T] {
	w := &Workers[T]{
		Work:    make(chan func() T),
		Results: make(chan T),
		wg:      sync.WaitGroup{},
	}

	for i := 0; i < numWorkers; i++ {
		w.wg.Add(1)
		go func() {
			for f := range w.Work {
				w.Results <- f()
			}
			w.wg.Done()
		}()
	}

	// Close the results channel when the work is done.
	go func() {
		w.wg.Wait()
		close(w.Results)
	}()

	return w
}
