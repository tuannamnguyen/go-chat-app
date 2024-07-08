package main

import (
	"context"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var wg sync.WaitGroup

func setupServer(hub *hub) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	e.GET("/chat/:chat_room/:user_name", hub.hubChatRoomHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hub := newHub(ctx, &wg)
	setupServer(hub)

	wg.Wait()
}
