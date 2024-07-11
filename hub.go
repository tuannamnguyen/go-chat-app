package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type hub struct {
	rooms map[string]*chatRoom
}

func newHub() *hub {
	return &hub{
		rooms: make(map[string]*chatRoom),
	}
}

func (h *hub) hubChatRoomHandler(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		chatRoom := c.Param("chat_room")
		userName := c.Param("user_name")

		room, ok := h.rooms[chatRoom]
		if !ok {
			user, err := newUser(userName, c.Response().Writer, c.Request())
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error creating user to new chat: %v", err))
			}

			room := h.addNewChatRoom(chatRoom)
			room.addUser(user)
			room.run(ctx)
		} else {
			if room.hasUser(userName) {
				log.Printf("%v already exists in room %v", userName, chatRoom)
			} else {
				user, err := newUser(userName, c.Response().Writer, c.Request())
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error creating user for chat: %v", err))
				} else {
					room.addUser(user)
					room.run(ctx)
				}
			}
		}
		return nil
	}
}

func (h *hub) addNewChatRoom(roomName string) *chatRoom {
	room := newChatRoom(roomName)
	h.rooms[roomName] = room
	return room
}
