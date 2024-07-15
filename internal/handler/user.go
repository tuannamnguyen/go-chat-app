package handler

import (
	"fmt"
	"net/http"

	"nhooyr.io/websocket"
)

type user struct {
	name      string
	conn      *websocket.Conn
	listening bool
}

func newUser(name string, w http.ResponseWriter, r *http.Request) (*user, error) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return nil, fmt.Errorf("error create new user: %v", err)
	}
	user := &user{
		name:      name,
		conn:      conn,
		listening: false,
	}

	return user, nil
}

func (u user) String() string {
	return fmt.Sprintf("user: %v", u.name)
}
