package models

import (
	"context"
	"fmt"
	"log"

	"nhooyr.io/websocket"
)

type ChatRoom struct {
	name         string
	users        []*user
	messages     chan message
	messagesRead []message
	addedUsers   chan *user
	dropUsers    chan *user
}

func NewChatRoom(roomName string) *ChatRoom {
	return &ChatRoom{
		name:       roomName,
		users:      []*user{},
		messages:   make(chan message, 100),
		dropUsers:  make(chan *user, 100),
		addedUsers: make(chan *user, 100),
	}
}

func (c *ChatRoom) AddUser(user *user) {
	c.addedUsers <- user
}

func (c *ChatRoom) Run(ctx context.Context) {
	go c.listen(ctx)
	go c.broadcast(ctx)
	go c.keepUserListUpdated(ctx)
}

func (c *ChatRoom) HasUser(userName string) bool {
	for _, user := range c.users {
		if user.name == userName {
			return true
		}
	}

	return false
}

func (c *ChatRoom) listen(ctx context.Context) {
	for {
		if len(c.users) > 0 {
			for _, user := range c.users {
				if !user.listening {
					user.listening = true
					go c.listenToUser(ctx, user)
				}
			}
		}
	}
}

func (c *ChatRoom) listenToUser(ctx context.Context, user *user) {
	for {
		log.Print("Listening to incoming messages")
		_, msg, err := user.conn.Read(ctx)
		if err != nil {
			log.Printf("error while listening to user messages: %v", err)
			log.Println(ctx.Err())
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

func (c *ChatRoom) broadcast(ctx context.Context) {
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
				log.Printf("error building message: %v, content: %s", err, bytes)
			} else {
				for _, user := range usersToSend {
					err := user.conn.Write(ctx, websocket.MessageText, bytes)
					if err != nil {
						log.Println(err)
					}
				}
				c.messagesRead = append(c.messagesRead, message)
			}
		case <-ctx.Done():
			break loop
		}
	}
}

func (c *ChatRoom) usersToSend(author *user) []*user {
	var result []*user
	for _, user := range c.users {
		if user != author {
			result = append(result, user)
		}
	}

	return result
}

func (c *ChatRoom) deleteUser(userToDelete *user) []*user {
	for i, user := range c.users {
		if userToDelete == user {
			result := append(c.users[:i], c.users[i+1:]...)
			return result
		}
	}

	return c.users
}

func (c *ChatRoom) keepUserListUpdated(ctx context.Context) {
loop:
	for {
		select {
		case user := <-c.addedUsers:
			c.users = append(c.users, user)
			c.broadcastMessage(ctx, []byte(fmt.Sprintf("%s joined %s\n", user.name, c.name)))
		case user := <-c.dropUsers:
			c.users = c.deleteUser(user)
			c.broadcastMessage(ctx, []byte(fmt.Sprintf("%s left %s\n", user.name, c.name)))
		case <-ctx.Done():
			break loop
		}

	}
}

func (c *ChatRoom) broadcastMessage(ctx context.Context, msg []byte) {
	for _, user := range c.users {
		err := user.conn.Write(ctx, websocket.MessageText, msg)
		if err != nil {
			log.Println(err)
		}
	}
}
