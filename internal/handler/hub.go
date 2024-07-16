package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/tuannamnguyen/go-chat-app/internal/models"

	"github.com/labstack/echo/v4"
)

type Hub struct {
	rooms map[string]*models.ChatRoom
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]*models.ChatRoom),
	}
}

func (h *Hub) HubChatRoomHandler(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		chatRoom := c.Param("chat_room")
		userName := c.Param("user_name")

		room, ok := h.rooms[chatRoom]
		if !ok {
			user, err := models.NewUser(userName, c.Response().Writer, c.Request())
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error creating user to new chat: %v", err))
			}

			room := h.AddNewChatRoom(chatRoom)
			room.AddUser(user)
			room.Run(ctx)
		} else {
			if room.HasUser(userName) {
				log.Printf("%v already exists in room %v", userName, chatRoom)
			} else {
				user, err := models.NewUser(userName, c.Response().Writer, c.Request())
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error creating user for chat: %v", err))
				} else {
					room.AddUser(user)
					room.Run(ctx)
				}
			}
		}
		return nil
	}
}

func (h *Hub) AddNewChatRoom(roomName string) *models.ChatRoom {
	room := models.NewChatRoom(roomName)
	h.rooms[roomName] = room
	return room
}
