package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func loginRouteHandler(c echo.Context) error {
	conf := &oauth2.Config{
		ClientID: os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL: os.Getenv("REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	url := conf.AuthCodeURL("state")
	return c.String(http.StatusOK, url)
}

func callbackRouteHandler(c echo.Context) error {
	return c.String(http.StatusOK, "abc")
}
