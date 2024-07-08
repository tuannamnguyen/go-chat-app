package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

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
	go func() {
		sigChan := make(chan os.Signal, 3)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-sigChan
		cancel()
	}()

	hub := newHub(ctx, &wg)
	setupServer(hub)

	wg.Wait()
}
