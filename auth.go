package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

type auth struct {
	config *oauth2.Config
	ctx    context.Context
}

func newAuth(ctx context.Context, config *oauth2.Config) *auth {
	return &auth{
		config: config,
		ctx:    ctx,
	}
}

func (a *auth) loginHandler(c echo.Context) error {
	url := a.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	return c.JSON(http.StatusOK, &apiResponse{
		Data: map[string]any{
			"auth_url": url,
		},
	})
}

func getUserInformation(ctx context.Context, client *http.Client) (*people.Person, error) {
	srv, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("error creating people api service: %v", err)
	}

	res, err := srv.People.Get("people/me").PersonFields("names").Do()
	if err != nil {
		return nil, fmt.Errorf("error getting call api get user info: %v", err)
	}

	return res, nil
}

func (a *auth) callbackHandler(c echo.Context) error {
	authCode := c.QueryParam("code")
	token, err := a.config.Exchange(a.ctx, authCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error when exchange auth code: %v", err))
	}

	apiClient := a.config.Client(c.Request().Context(), token)
	userInfo, err := getUserInformation(a.ctx, apiClient)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("error when get user information: %v", err))
	}

	return c.JSON(http.StatusOK, &apiResponse{Data: map[string]any{"user_data": userInfo}})
}
