package task

// This type implements a task abstraction loosely based on .NET Tasks.
type Task struct {
	f       func()
	awaiter chan struct{}
}

// Creates a new task.
func New(f func()) *Task {
	return &Task{
		f:       f,
		awaiter: make(chan struct{}),
	}
}

// Executes the task.
func (t *Task) Start() {
	go func() {
		t.f()
		t.awaiter <- struct{}{}
	}()
}

// Waits for the task to complete.
func (t *Task) Wait() {
	<-t.awaiter
}

// Creates a new task with a single input argument.
func NewWithInput[T any](f func(T), input T) *Task {
	fun := func() {
		f(input)
	}
	return New(fun)
}

// Task struct that stores a single result value.
type TaskWithResult[T any] struct {
	Task
	result T
}

// Creates a new task that stores a single result value.
func NewWithResult[T any](f func() T) *TaskWithResult[T] {
	t := &TaskWithResult[T]{
		Task: *New(nil),
	}
	t.f = func() {
		t.result = f()
	}
	return t
}

// Returns the result value after waiting for the task to finish.
func (t *TaskWithResult[T]) GetResult() T {
	t.Wait()
	return t.result
}
