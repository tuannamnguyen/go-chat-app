package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dotenv-org/godotenvvault"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tuannamnguyen/go-chat-app/internal/handler"
	"github.com/tuannamnguyen/go-chat-app/internal/repository"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/people/v1"
)

func setupServer(ctx context.Context, e *echo.Echo, hub *handler.HubService, auth *handler.AuthService) {
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	e.GET("/chat/:chat_room/:user_name", hub.HubChatRoomHandler(ctx))

	e.GET("/auth/login", auth.LoginHandler)
	e.GET("/auth/callback", auth.CallbackHandler)
	e.GET("/auth/user/:user_id", auth.GetUserName)

	if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
		e.Logger.Fatal("shutting down the server")
	}
}

func main() {
	//setup .env
	err := godotenvvault.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

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
	repository := repository.NewAuthRepository()

	//setup server
	e := echo.New()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	defer stop()

	hub := handler.NewHubService()
	auth := handler.NewAuthService(config, repository)
	go setupServer(ctx, e, hub, auth)

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
