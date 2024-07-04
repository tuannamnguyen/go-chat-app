package main

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func generateStateOauthCookie() string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	return state
}

func loginRouteHandler(c echo.Context) error {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	url := conf.AuthCodeURL(generateStateOauthCookie(), oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	return c.JSON(http.StatusOK, &apiResponse{Data: url})
}

func callbackRouteHandler(c echo.Context) error {
	return c.String(http.StatusOK, "abcd")
}
