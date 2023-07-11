package actor

import (
	"sync"
)

type Actor[T any] struct {
	messages chan T
	handler  func(T)
	wg       sync.WaitGroup
}

// Creates a new actor with the given handler function.
func New[T any](handler func(T)) *Actor[T] {
	a := &Actor[T]{
		messages: make(chan T),
		handler:  handler,
	}
	go func() {
		for m := range a.messages {
			a.handler(m)
			a.wg.Done()
		}
	}()
	return a
}

// Sends a message to the actor.
func (a *Actor[T]) Send(m T) {
	a.wg.Add(1)
	go func() {
		a.messages <- m
	}()
}

// Waits for all messages to be finished and closes the channel.
func (a *Actor[T]) Stop() {
	a.wg.Wait()
	close(a.messages)
}
