package main

import (
	"fmt"
	"os"

	"github.com/turmaxx/concurrency-abstractions-in-go/actor"
	"github.com/turmaxx/concurrency-abstractions-in-go/events"
	"github.com/turmaxx/concurrency-abstractions-in-go/task"
	"github.com/turmaxx/concurrency-abstractions-in-go/workers"
)

func main() {
	switch os.Args[1] {
	// Helloworld example
	case "helloworld":
		helloworld()
		// Task examples
	case "task":
		t := task.New(func() {
			fmt.Println("Hello, World!")
		})
		t.Start()
		t.Wait()
	case "task_input":
		greeting := "Hello from the closure!"
		t := task.New(func() {
			fmt.Println(greeting)
		})
		t.Start()
		t.Wait()
	case "task_output":
		var greeting string
		t := task.New(func() {
			greeting = "Hello from the closure!"
		})
		t.Start()
		t.Wait()
		fmt.Println(greeting)
	case "task_input_generic":
		t := task.NewWithInput(func(i string) {
			fmt.Println(i)
		}, "Hello with generics!")
		t.Start()
		t.Wait()
	case "task_output_generic":
		t := task.NewWithResult(func() string {
			return "Hello with generic output!"
		})
		t.Start()
		fmt.Println(t.GetResult())
	case "task_nested":
		t := task.New(func() {
			fmt.Println("Hello from the first task!")
			t := task.New(func() {
				fmt.Println("Hello from the second task!")
			})
			t.Start()
			t.Wait()
		})
		t.Start()
		t.Wait()
	// Worker examples
	case "worker":
		w := workers.New[string](2)
		go func() {
			for i := 0; i < 10; i++ {
				i := i
				w.Work <- func() string {
					return fmt.Sprintf("%d", i)
				}
			}
			close(w.Work)
		}()
		for r := range w.Results {
			fmt.Println(r)
		}
		// Actor examples
	case "actor":
		a := actor.New(func(s string) {
			fmt.Println(s)
		})
		a.Send("Hello, World!")
		a.Send("Hello again, World!")
		a.Stop()
	case "actor_printer":
		p := actor.NewPrinter()
		p.Print("Hello, World!")
		p.Print("Hello again, World!")
		p.Stop()
	case "actor_chatroom":
		// create chat room
		chatRoom := actor.NewChatRoom()

		// create clients
		alice := actor.NewClient("Alice", chatRoom)
		bob := actor.NewClient("Bob", chatRoom)

		// send messages
		alice.Send("Hello, Bob!")
		bob.Send("Hello, Alice!")

		// stop actors and wait for them to finish
		alice.Stop()
		bob.Stop()
		chatRoom.Stop()
	// Events examples
	case "pubsub":
		p := events.NewPublisher[string]()
		done1 := make(chan struct{})
		done2 := make(chan struct{})
		s1 := events.NewSubscriber(p, func(s string) {
			fmt.Println("Subscriber 1:", s)
			done1 <- struct{}{}
		})
		events.NewSubscriber(p, func(s string) {
			fmt.Println("Subscriber 2:", s)
			done2 <- struct{}{}
		})
		p.Publish("Hello, World!")
		<-done1
		<-done2
		p.Wait()
		s1.Unsubscribe()
		p.Publish("Hello, again World!")
		<-done2
		p.Wait()
	default:
		fmt.Println("Unknown example")
	}
}
