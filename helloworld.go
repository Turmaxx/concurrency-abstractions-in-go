package main

import "fmt"

func helloworld() {
	c := make(chan struct{})
	go func() {
		fmt.Println("Hello, World!")
		c <- struct{}{}
	}()
	<-c
}
