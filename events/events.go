package events

import "sync"

type Publisher[T any] struct {
	sync.Mutex
	subscribers map[int]chan T
	lastID      int
	wg          sync.WaitGroup
}

func NewPublisher[T any]() *Publisher[T] {
	return &Publisher[T]{
		subscribers: make(map[int]chan T),
	}
}

// Publish a message to all subscribers.
func (p *Publisher[T]) Publish(m T) {
	p.Lock()
	defer p.Unlock()

	for _, c := range p.subscribers {
		p.wg.Add(1)
		go func(c chan T) {
			defer p.wg.Done()
			c <- m
		}(c)
	}
}

// Close all subscriber channels
func (p *Publisher[T]) Stop() {
	p.Lock()
	defer p.Unlock()

	for _, c := range p.subscribers {
		close(c)
	}
}

func (p *Publisher[T]) Wait() {
	p.wg.Wait()
}

func (p *Publisher[T]) subscribe(s *Subscriber[T]) int {
	p.Lock()
	defer p.Unlock()

	id := p.lastID + 1
	p.subscribers[id] = s.channel
	p.lastID = id
	return id
}

func (p *Publisher[T]) unsubscribe(id int) {
	p.Lock()
	defer p.Unlock()

	delete(p.subscribers, id)
}

type Subscriber[T any] struct {
	callback  func(T)
	id        int
	channel   chan T
	publisher *Publisher[T]
}

func NewSubscriber[T any](p *Publisher[T], cb func(T)) *Subscriber[T] {
	s := &Subscriber[T]{
		publisher: p,
		channel:   make(chan T),
		callback:  cb,
	}
	s.id = p.subscribe(s)
	go func() {
		for m := range s.channel {
			s.callback(m)
		}
	}()
	return s
}

func (s *Subscriber[T]) Unsubscribe() {
	s.publisher.unsubscribe(s.id)
	close(s.channel)
}
