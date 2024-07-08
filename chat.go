package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"nhooyr.io/websocket"
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

func newChatRoom(roomName string, ctx context.Context, wg *sync.WaitGroup) *chatRoom {
	return &chatRoom{
		name:       roomName,
		users:      []*user{},
		messages:   make(chan message, 100),
		dropUsers:  make(chan *user, 100),
		addedUsers: make(chan *user, 100),
		ctx:        ctx,
		wg:         wg,
	}
}

func (c *chatRoom) addUser(user *user) {
	c.addedUsers <- user
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
	c.wg.Add(1)
	for {
		log.Print("Listening to incoming messages")
		_, msg, err := user.conn.Read(c.ctx)
		if err != nil {
			log.Printf("error while listening to user messages: %v", err)
			c.dropUsers <- user
			break

		} else {
			c.messages <- message{
				bytes:  msg,
				author: user,
			}
		}
	}
}

func (c *chatRoom) broadcast() {
	c.wg.Add(1)
	log.Println("broadcasting messages")

loop:
	for {
		select {
		case message := <-c.messages:
			log.Printf("received message from: %v", message.author.name)

			usersToSend := c.usersToSend(message.author)
			log.Printf("broadcasting message to: %v", usersToSend)
			bytes, err := message.prepareMsg()
			if err != nil {
				for _, user := range usersToSend {
					user.conn.Write(c.ctx, websocket.MessageText, bytes)
				}
				c.messagesRead = append(c.messagesRead, message)
			} else {
				log.Printf("error building message: %v, content: %s", err, bytes)
			}
		case <-c.ctx.Done():
			break loop
		}
	}

	c.wg.Done()
}

func (c *chatRoom) usersToSend(author *user) []*user {
	var result []*user
	for _, user := range c.users {
		if user != author {
			result = append(result, user)
		}
	}

	return result
}

func (c *chatRoom) deleteUser(userToDelete *user) []*user {
	for i, user := range c.users {
		if userToDelete == user {
			result := append(c.users[:i], c.users[i+1:]...)
			return result
		}
	}

	return c.users
}

func (c *chatRoom) keepUserListUpdated() {
loop:
	for {
		select {
		case user := <-c.addedUsers:
			c.users = append(c.users, user)
			c.broadcastMessage([]byte(fmt.Sprintf("%s joined %s\n", user.name, c.name)))
		case user := <-c.dropUsers:
			c.users = c.deleteUser(user)
			c.broadcastMessage([]byte(fmt.Sprintf("%s left %s\n", user.name, c.name)))
		case <-c.ctx.Done():
			break loop
		}

	}
}

func (c *chatRoom) broadcastMessage(msg []byte) {
	for _, user := range c.users {
		user.conn.Write(c.ctx, websocket.MessageText, msg)
	}
}
