package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/tuannamnguyen/go-chat-app/internal/models"

	"github.com/labstack/echo/v4"
)

type HubService struct {
	hub *models.Hub
}

func NewHubService() *HubService {
	hub := models.NewHub()

	return &HubService{hub}
}

func (h *HubService) HubChatRoomHandler(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		chatRoom := c.Param("chat_room")
		userName := c.Param("user_name")

		room, ok := h.hub.Rooms[chatRoom]
		if !ok {
			user, err := models.NewUser(userName, c.Response().Writer, c.Request())
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("error creating user to new chat: %v", err))
			}

			room := h.hub.AddNewChatRoom(chatRoom)
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
