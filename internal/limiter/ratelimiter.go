package limiter

import (
	"sync"
)

type Args interface {
	string | int | []string | []int
}

type Gatherer[T any] func(int) T

type Limiter[T any] struct {
	dataChannel  chan T
	function     Gatherer[T]
	jobChannel   chan int
	limit        int
	numberOfJobs int
	waitGroup    *sync.WaitGroup
}

/*
Limiter takes a `limit`, a `numberOfJobs` and a worker func that gathers data.

The `limit` determines the max number of goroutines to be spawned.

The `numberOfJobs` determines how many times the worker func will be invoked.

The Spawn.Run() chained method will return a data channel of the return type of the worker func.

Usage:

	limiter := New(maxGoroutines, numberOfJobs, myWorkerFunc)
	results := limiter.Spawn().Run()

	for result := range results {
		doTheNeedful(result)
	}
*/
func New[T any](limit, numberOfJobs int, fn Gatherer[T]) *Limiter[T] {
	return &Limiter[T]{
		dataChannel:  make(chan T, numberOfJobs),
		function:     fn,
		jobChannel:   make(chan int, limit),
		limit:        limit,
		numberOfJobs: numberOfJobs,
		waitGroup:    &sync.WaitGroup{},
	}
}

func (l *Limiter[T]) Gather(j int) {
	defer l.waitGroup.Done()
	for i := range l.jobChannel {
		data := l.function(i)
		l.dataChannel <- data
	}
}

func (l *Limiter[T]) Spawn() *Limiter[T] {
	for i := 0; i < l.limit; i++ {
		l.waitGroup.Add(1)
		go l.Gather(i)
	}
	return l
}

func (l *Limiter[T]) Run() chan T {
	for i := 0; i < l.numberOfJobs; i++ {
		l.jobChannel <- i
	}
	close(l.jobChannel)
	l.waitGroup.Wait()
	close(l.dataChannel)
	return l.dataChannel
}
