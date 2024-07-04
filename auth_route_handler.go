package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type authController struct {
	authConfig *oauth2.Config
}

func newAuthController() *authController {
	return &authController{authConfig: &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}}
}

func generateStateOauthCookie() string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	return state
}

func (a *authController) loginRouteHandler(c echo.Context) error {
	url := a.authConfig.AuthCodeURL(generateStateOauthCookie(), oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	return c.String(http.StatusOK, url)
}

func (a *authController) callbackRouteHandler(c echo.Context) error {
	authCode := c.FormValue("code")
	token, err := exchangeCodeForToken(authCode, a.authConfig)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, &apiResponse{Data: token.AccessToken})
}

func exchangeCodeForToken(authCode string, config *oauth2.Config) (*oauth2.Token, error) {
	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
	}

	return token, nil
}
