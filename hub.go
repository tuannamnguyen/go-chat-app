package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

type hub struct {
	rooms map[string]*chatRoom
	ctx   context.Context
	wg    *sync.WaitGroup
}

func newHub(ctx context.Context, wg *sync.WaitGroup) *hub {
	return &hub{
		rooms: make(map[string]*chatRoom),
		ctx:   ctx,
		wg:    wg,
	}
}

func (h *hub) chatRoom(c echo.Context) error {
	chatRoom := c.Param("chat_room")
	userName := c.Param("user_name")

	room, ok := h.rooms[chatRoom]
	if !ok {
		room := h.addChatRoom(chatRoom)
		user, err := newUser(userName, c.Response().Writer, c.Request())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error creating user to new chat: %v", err))
		}
		room.addUser(user)
		room.run()
	} else {
		if room.hasUser(userName) {
			return c.JSON(http.StatusOK, &apiResponse{Data: "User already exists in chat room"})
		} else {
			user, err := newUser(userName, c.Response().Writer, c.Request())
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error creating user for chat: %v", err))
			} else {
				room.addUser(user)
				room.run()
			}
		}
	}
	return nil
}

func (h *hub) addChatRoom(roomName string) *chatRoom {
	room := newChatRoom(roomName, h.ctx, h.wg)
	h.rooms[roomName] = room
	return room
}
