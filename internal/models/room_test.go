package models

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"nhooyr.io/websocket"
)

func TestChatRoom_listenToUser(t *testing.T) {
	type fields struct {
		name         string
		users        []*user
		messages     chan message
		messagesRead []message
		addedUsers   chan *user
		dropUsers    chan *user
	}
	tests := []struct {
		name     string
		fields   fields
		wantMsg  message
		wantDrop bool
		send     func(*testing.T, *websocket.Conn)
	}{
		{
			name: "Successfully read message",
			fields: fields{
				name:      "TestRoom",
				messages:  make(chan message, 1),
				dropUsers: make(chan *user, 1),
			},
			wantMsg: message{
				bytes: []byte("Hello, World!"),
			},
			send: func(t *testing.T, conn *websocket.Conn) {
				err := conn.Write(context.Background(), websocket.MessageText, []byte("Hello, World!"))
				if err != nil {
					t.Fatalf("Failed to send message: %v", err)
				}
			},
		},
		{
			name: "Connection closed",
			fields: fields{
				name:      "TestRoom",
				messages:  make(chan message, 1),
				dropUsers: make(chan *user, 1),
			},
			wantDrop: true,
			send: func(t *testing.T, conn *websocket.Conn) {
				err := conn.Close(websocket.StatusNormalClosure, "")
				if err != nil {
					t.Fatalf("Failed to close connection: %v", err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				conn, err := websocket.Accept(w, r, nil)
				if err != nil {
					t.Fatalf("Failed to accept connection: %v", err)
				}
				if tt.send != nil {
					tt.send(t, conn)
				}
			}))
			defer s.Close()

			// Create a client connection
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			conn, _, err := websocket.Dial(ctx, "ws"+s.URL[4:], nil)
			if err != nil {
				t.Fatalf("Failed to dial: %v", err)
			}

			c := &ChatRoom{
				name:         tt.fields.name,
				users:        tt.fields.users,
				messages:     tt.fields.messages,
				messagesRead: tt.fields.messagesRead,
				addedUsers:   tt.fields.addedUsers,
				dropUsers:    tt.fields.dropUsers,
			}

			u := &user{
				name: "TestUser",
				conn: conn,
			}

			go c.listenToUser(ctx, u)

			if tt.wantMsg.bytes != nil {
				select {
				case msg := <-c.messages:
					if string(msg.bytes) != string(tt.wantMsg.bytes) {
						t.Errorf("Expected message %s, got %s", tt.wantMsg.bytes, msg.bytes)
					}
				case <-time.After(time.Second):
					t.Error("Timeout waiting for message")
				}
			}

			if tt.wantDrop {
				select {
				case droppedUser := <-c.dropUsers:
					if droppedUser != u {
						t.Errorf("Expected dropped user %v, got %v", u, droppedUser)
					}
				case <-time.After(time.Second):
					t.Error("Timeout waiting for user to be dropped")
				}
			}
		})
	}
}
