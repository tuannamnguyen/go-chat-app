package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func setupServer() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
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

	setupServer()

}
