package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type auth struct {
	config *oauth2.Config
}

func newAuth(config *oauth2.Config) *auth {
	return &auth{
		config: config,
	}
}

func (a *auth) loginHandler(c echo.Context) error {
	url := a.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	return c.JSON(http.StatusOK, &apiResponse{Data: url})
}
