package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func setupServer() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	authController := newAuthController()
	authGroup := e.Group("/auth")
	authGroup.GET("/login", authController.loginRouteHandler)
	authGroup.GET("/callback", authController.callbackRouteHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("can't load .env file")
	}

	key := "secret-session-key"
	store := sessions.NewCookieStore([]byte(key))
	gothic.Store = store

	goth.UseProviders(google.New(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"), os.Getenv("REDIRECT_URL")))

	setupServer()

}
