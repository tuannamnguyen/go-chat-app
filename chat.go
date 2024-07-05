package main

import (
	"context"
	"sync"
)

type chatRoom struct {
	name         string
	users        []*user
	messages     chan message
	messagesRead []message
	addedUsers   chan *user
	dropUsers    chan *user
	ctx          context.Context
	wg           *sync.WaitGroup
}

func (c *chatRoom) addUser(user *user) {
	c.users = append(c.users, user)
}

func (c *chatRoom) run() {
	go c.listen()
	go c.broadcast()
	go c.keepUserListUpdated()
}

func (c *chatRoom) hasUser(userName string) bool {
	for _, user := range c.users {
		if user.name == userName {
			return true
		}
	}

	return false
}

func (c *chatRoom) listen() {
	for {
		if len(c.users) > 0 {
			for _, user := range c.users {
				if !user.listening {
					user.listening = true
					go c.listenToUser(user)
				}
			}
		}
	}
}

func (c *chatRoom) listenToUser(user *user) {
	for {
		_, msg, err := user.conn.Read(c.ctx)
		if err == nil {
			c.messages <- message{
				bytes: msg,
				author: user,
			}
		} 
	}
} 

func (c *chatRoom) broadcast() {}

func (c *chatRoom) keepUserListUpdated() {}
