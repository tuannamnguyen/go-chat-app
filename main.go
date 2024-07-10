package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	prettylogger "github.com/rdbell/echo-pretty-logger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/people/v1"
)

var wg sync.WaitGroup

func setupServer(e *echo.Echo, hub *hub, auth *auth) {
	e.Use(prettylogger.Logger)
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	e.GET("/chat/:chat_room/:user_name", hub.hubChatRoomHandler)

	e.GET("/auth/login", auth.loginHandler)
	e.GET("/auth/callback", auth.callbackHandler)
	e.GET("/auth/user/:user_id", auth.getUserName)

	if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
		e.Logger.Fatal("shutting down the server")
	}
}

func main() {
	//setup oauth2
	config := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes: []string{
			people.UserinfoProfileScope,
			people.UserinfoEmailScope,
		},
		Endpoint: google.Endpoint,
	}

	//setup redis
	redisHandler := newRedisHandler()

	//setup server
	e := echo.New()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	defer stop()

	hub := newHub(&wg)
	auth := newAuth(config, redisHandler)
	go setupServer(e, hub, auth)

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	wg.Wait()
}
