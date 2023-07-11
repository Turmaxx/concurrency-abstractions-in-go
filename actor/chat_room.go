package actor

import (
	"fmt"
)

type ChatRoom struct {
	*Actor[message]
}

// Creates a new chat room that prints messages that come in.
func NewChatRoom() *ChatRoom {
	c := &ChatRoom{}
	c.Actor = New(func(m message) {
		fmt.Printf("%s: \t %s\n", m.sender.name, m.text)
	})
	return c
}

type Client struct {
	*Actor[string]
	name string
}

// Creates a new client that sends messages to the given chat room.
func NewClient(name string, room *ChatRoom) *Client {
	c := &Client{
		name: name,
	}
	c.Actor = New(func(t string) {
		m := message{
			sender: c,
			text:   t,
		}
		room.Send(m)
	})
	return c
}

type message struct {
	sender *Client
	text   string
}
